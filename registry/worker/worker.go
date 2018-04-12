package registryworker

import (
	"log"

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

func Notify(ch <-chan *Event) <-chan error {
	out := make(chan error)
	log.Printf("registry worker to notify start")
	go func() {
		for e := range ch {
			inwh, err := slack.NewRegistryIncomingWebhook(e.Event, e.Err)
			if err != nil {
				log.Printf("[error] %v", err)
				continue
			}
			if err := notifier(inwh); err != nil {
				log.Printf("[error] %v", err)
				continue
			}
			out <- nil
		}
	}()
	log.Printf("registry worker to notify stopped")
	return out
}