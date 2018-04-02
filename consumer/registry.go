package consumer

import (
	"fmt"
	"log"
	"strings"

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
	prefix := strings.ToUpper(viper.GetString("prefix"))
	f := func(a string, b interface{}) string { return fmt.Sprintf("%s%s=%v", prefix, a, b) }
	return []string{
		f("EVENT_ID", e.ID),
		f("REPOSITORY", e.Target.Repository),
		f("DIGEST", e.Target.Digest),
		f("URL", e.Target.URL),
		f("TAG", e.Target.Tag),
		f("REQUEST_ID", e.Request.ID),
		f("ADDR", e.Request.Addr),
		f("HOST", e.Request.Host),
		f("USER_AGENT", e.Request.UserAgent),
	}
}
