package consumer

import (
	"log"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque"
	"github.com/tmtk75/weque/registry"
	"github.com/tmtk75/weque/repository"
	"github.com/tmtk75/weque/slack"
)

const (
	KeyHandlerScriptRegistry   = "handlers.registry"
	KeyHandlerScriptRepository = "handlers.repository"
)

func StartRepositoryConsumer(events <-chan *repository.Context, out chan<- *Event) {
	log.Printf("repository consumer started")
	for c := range events {
		w := c.Webhook
		log.Printf("started. delivery: %v", w.Delivery)

		rh := viper.GetString(KeyHandlerScriptRepository)
		err := weque.Run(w.Env(), ".", rh)
		if err != nil {
			log.Printf("[error] %v", err)
		}
		out <- &Event{Context: c, Err: err}
		log.Printf("finished. id: %v", w.Delivery)
	}
}

var notifier = slack.Notify

func StartRegistryConsumer(events <-chan *registry.Webhook, out chan<- *Event) {
	log.Printf("registry consumer started")
	for w := range events {
		e := w.Events[0]
		log.Printf("started. id: %v", e.ID)

		rh := viper.GetString(KeyHandlerScriptRegistry)
		err := weque.Run(e.Env(), ".", rh)
		if err != nil {
			log.Printf("[error] %v", err)
		}
		out <- &Event{Event: &e, Err: err}
		log.Printf("finished. id: %v", e.ID)
	}
}

type Event struct {
	Event   *registry.Event
	Context *repository.Context
	Err     error
}
