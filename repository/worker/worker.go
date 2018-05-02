package repositoryworker

import (
	"log"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque"
	"github.com/tmtk75/weque/repository"
	"github.com/tmtk75/weque/slack"
)

const (
	KeyHandlerScriptRepository = "handlers.repository"
)

type Context struct {
	*repository.Context
	Err error
}

func Run(events <-chan *repository.Context) <-chan *Context {
	out := make(chan *Context)
	go func() {
		for c := range events {
			w := c.Webhook
			log.Printf("started. delivery: %v", w.Delivery)

			rh := viper.GetString(KeyHandlerScriptRepository)
			err := weque.Run(w.Env(), ".", rh)
			if err != nil {
				log.Printf("[error] %v", err)
			}
			out <- &Context{Context: c, Err: err}
			log.Printf("finished. id: %v", w.Delivery)
		}
		close(out)
		log.Printf("repository worker to run stopped")
	}()
	log.Printf("repository worker to run started")
	return out
}

var notifier = slack.Notify

func Notify(ch <-chan *Context) <-chan error {
	out := make(chan error)
	go func() {
		for e := range ch {
			inwh, err := slack.NewIncomingWebhookRepository(e.Context.Webhook, e.Context.WebhookProvider, e.Err)
			if err != nil {
				//log.Printf("[error] failed to build incoming webhook: %v", err)
				out <- err
				continue
			}
			if err := notifier(inwh); err != nil {
				//log.Printf("[error] failed to notify: %v", err)
				out <- err
				continue
			}
			log.Printf("notified: %v", e.Context.Webhook.Delivery)
			out <- nil
		}
		close(out)
		log.Printf("repository worker to notify stopped")
	}()
	log.Printf("repository worker to notify started")
	return out
}
