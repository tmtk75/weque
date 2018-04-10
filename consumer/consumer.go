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

func Start(repo <-chan *repository.Context, reg <-chan *registry.Webhook) {
	go StartRepositoryConsumer(repo)
	go StartRegistryConsumer(reg)
}

func StartRepositoryConsumer(events <-chan *repository.Context) {
	log.Printf("repository consumer started")
	for {
		ctx := <-events
		wb := ctx.Webhook
		log.Printf("started. delivery: %v", wb.Delivery)

		rh := viper.GetString(KeyHandlerScriptRepository)

		err := weque.Run(wb.Env(), ".", rh)
		if err != nil {
			log.Printf("[error] %v", err)
		}
		log.Printf("finished. delivery: %v", wb.Delivery)

		if err := slack.Notify(ctx, err); err != nil {
			log.Printf("[error] %v", err)
		}
	}
}

func StartRegistryConsumer(events <-chan *registry.Webhook) {
	log.Printf("registry consumer started")
	for {
		wh := <-events
		e := &wh.Events[0] // 1 event is ensured
		log.Printf("started. id: %v", e.ID)

		rh := viper.GetString(KeyHandlerScriptRegistry)

		err := weque.Run(e.Env(), ".", rh)
		if err != nil {
			log.Printf("[error] %v", err)
		}
		log.Printf("finished. id: %v", e.ID)

		if err := slack.NotifyRegistry(e, err); err != nil {
			log.Printf("[error] %v", err)
		}
	}
}
