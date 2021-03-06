package registry

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque"
)

const (
	KeySecretToken  = weque.KeySecretToken
	KeyInsecureMode = weque.KeyInsecureMode
)

func NewHandler(events chan<- *Webhook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !viper.GetBool(KeyInsecureMode) {
			if err := Verify(r); err != nil {
				weque.SendError(w, 401, fmt.Sprintf("failed to verify: %v", err))
				return
			}
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read: %v", err)
			return
		}
		//log.Printf("payload: %v", string(b))

		wh, err := Unmarshal(b)
		if err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}

		//payload, err := json.Marshal(wh)
		//if err != nil {
		//	log.Printf("failed to marshal for debug log: %v", err) // continue to process because just debug log
		//}
		//log.Printf("payload: %v", string(payload))

		if len(wh.Events) != 1 {
			if len(wh.Events) == 0 {
				log.Printf("no any events")
				weque.SendError(w, 400, "no any events")
				return
			}
			ids := make([]string, len(wh.Events))
			for i, e := range wh.Events {
				ids[i] = e.ID
			}
			rb := strings.Join(ids, ",")
			log.Printf("doesn't support multiple events: %v", rb)
			weque.SendError(w, 400, fmt.Sprintf("unsupported multiple events: %v", rb))
			return
		}

		e := wh.Events[0]
		log.Printf("repository: %v, tag: %v, media-type: %v, action: %v", e.Target.Repository, e.Target.Tag, e.Target.MediaType, e.Action)

		go (func() {
			if !(e.Action == "push" && e.Target.MediaType == "application/vnd.docker.distribution.manifest.v2+json") {
				log.Printf("ignored action. id: %v, media-type: %v, action: %v", e.ID, e.Target.MediaType, e.Action)
				return
			}
			events <- wh
			log.Printf("queued. id: %v", e.ID)
		})()

		log.Printf("ok: %v", e.ID)
		w.Write([]byte(e.ID))
	}
}

func Verify(r *http.Request) error {
	token := viper.GetString(KeySecretToken)
	secret := r.Header.Get("x-weque-secret")
	if token != secret {
		return fmt.Errorf("the given secret token didn't match")
	}
	return nil
}
