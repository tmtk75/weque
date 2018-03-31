package registry

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	KeyHandlerScript = "handlers.registry"
)

func NewHandler(events chan<- *Webhook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read: %v", err)
			return
		}

		wh, err := Unmarshal(b)
		if err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}

		if len(wh.Events) != 1 {
			ids := make([]string, len(wh.Events))
			for i, e := range wh.Events {
				ids[i] = e.ID
			}
			rb := strings.Join(ids, ",")
			log.Printf("doesn't support multiple events: %v", rb)
			w.Write([]byte(fmt.Sprintf("unsupported multiple events: %v", rb)))
			return
		}

		e := wh.Events[0]

		go (func() {
			events <- wh
			log.Printf("queued. id: %v", e.ID)
		})()

		log.Printf("ok: %v", e.ID)
		w.Write([]byte(e.ID))
	}
}
