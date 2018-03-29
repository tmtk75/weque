package consumer

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
