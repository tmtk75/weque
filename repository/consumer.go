package repository

import (
	"log"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque"
)

const (
	KeyHandlerScript = "handlers.repository"
)

func StartConsumer(events <-chan *Webhook) {
	log.Printf("repository consumer started")
	for {
		wb := <-events
		log.Printf("started. delivery: %v", wb.Delivery)

		rh := viper.GetString(KeyHandlerScript)

		err := weque.Run(wb.Env(), ".", rh)
		if err != nil {
			log.Printf("[error] %v", err)
		}
		log.Printf("finished. delivery: %v", wb.Delivery)
	}
}
