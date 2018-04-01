package consumer

import (
	"log"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque"
	"github.com/tmtk75/weque/registry"
	"github.com/tmtk75/weque/slack"
)

const (
	KeyHandlerScriptRegistry = "handlers.registry"
)

func StartRegistryConsumer(events <-chan *registry.Webhook) {
	log.Printf("registry consumer started")
	for {
		wh := <-events
		e := &wh.Events[0] // 1 event is ensured
		log.Printf("started. id: %v", e.ID)

		rh := viper.GetString(KeyHandlerScriptRegistry)

		err := weque.Run(RegistryEnv(e), ".", rh)
		if err != nil {
			log.Printf("[error] %v", err)
		}
		log.Printf("finished. id: %v", e.ID)

		if err := slack.NotifyRegistry(e, err); err != nil {
			log.Printf("[error] %v", err)
		}
	}
}

func RegistryEnv(e *registry.Event) []string {
	return make([]string, 9)
}
