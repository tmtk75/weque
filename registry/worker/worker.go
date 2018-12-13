package registryworker

import (
	"log"
	"regexp"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque"
	"github.com/tmtk75/weque/registry"
	"github.com/tmtk75/weque/slack"
)

const KeyHandlerScriptRegistry = "handlers.registry"

type Event struct {
	*registry.Event
	Err error
}

func Run(events <-chan *registry.Webhook) <-chan *Event {
	out := make(chan *Event)
	go func() {
		for w := range events {
			e := w.Events[0] // ensured
			log.Printf("started. id: %v", e.ID)

			rh := viper.GetString(KeyHandlerScriptRegistry)
			err := weque.Run(e.Env(), ".", rh)
			if err != nil {
				log.Printf("[error] %v", err)
			}
			out <- &Event{Event: &e, Err: err}
			log.Printf("finished. id: %v", e.ID)
		}
		close(out)
		log.Printf("registry worker to run stopped")
	}()
	log.Printf("registry worker to run started")
	return out
}

var notifier = slack.Notify

const KeySlackNotificationRegistryExclude = "notification.slack.registry.exclude"

func Notify(ch <-chan *Event) <-chan error {
	out := make(chan error)
	go func() {
		for e := range ch {
			inwh, err := slack.NewIncomingWebhookRegistry(e.Event, e.Err)
			if err != nil {
				//log.Printf("[error] %v", err)
				out <- err
				continue
			}

			exclude := viper.GetString(KeySlackNotificationRegistryExclude)
			log.Printf("may be excluded by %v", exclude)
			if exclude != "" && Exclude(e.Event, exclude) {
				log.Printf("not notified to slack for %v of %v", ExcludeTarget(e.Event), e.Event.ID)
				continue
			}

			if err := notifier(inwh); err != nil {
				//log.Printf("[error] %v", err)
				out <- err
				continue
			}
			log.Printf("notified. event: %v", e.Event)
			out <- nil
		}
		close(out)
		log.Printf("registry worker to notify stopped")
	}()
	log.Printf("registry worker to notify started")
	return out
}

func ExcludeTarget(e *registry.Event) string {
	if e == nil {
		return ""
	}
	return e.Target.Repository + ":" + e.Target.Tag
}

func Exclude(e *registry.Event, restr string) bool {
	if e == nil {
		return true
	}
	re, err := regexp.Compile(restr)
	if err != nil {
		return true
	}
	return re.Match([]byte(ExcludeTarget(e)))
}
