package repository

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque"
)

type Handler interface {
	/*
	 * Return the requiest ID string
	 */
	RequestID(r *http.Request) string

	/*
	 * Return nil if the request is acceptable.
	 */
	Verify(r *http.Request, body []byte) error

	/*
	 * Returns true if the request is just a ping but not actual webhook.
	 */
	IsPing(r *http.Request, body []byte) bool

	/*
	 * Returns a Webhook for the given response body.
	 */
	Unmarshal(r *http.Request, body []byte) (*Webhook, error)

	/** */
	WebhookProvider() WebhookProvider
}

const (
	KeyInsecureMode = weque.KeyInsecureMode
	KeySecretToken  = weque.KeySecretToken
)

func Shorten(s string, n int) string {
	if len(s) > n {
		return s[0:n] + "..."
	}
	return s
}

func NewHandler(h Handler, events chan<- *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rid := h.RequestID(r)

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			weque.SendError(w, 500, fmt.Sprintf("failed to read: %v", err))
			return
		}
		log.Printf("[debug] read %vbytes for %v", len(b), rid)

		if !viper.GetBool(KeyInsecureMode) {
			if err := h.Verify(r, b); err != nil {
				log.Printf("failed to verify for %v: %v", rid, Shorten(string(b), 32))
				weque.SendError(w, 401, fmt.Sprintf("failed to verify: %v", err))
				return
			}
			log.Printf("[debug] verified %v", rid)
		} else {
			log.Printf("skip to verify because of insecure mode")
		}

		if h.IsPing(r, b) {
			log.Printf("request is ping: %v", rid)
			msg := fmt.Sprintf("Received ping: %v", r.RequestURI)
			log.Print(msg)
			w.WriteHeader(200)
			w.Write([]byte(msg))
			return
		}

		body, err := h.Unmarshal(r, b)
		if err != nil {
			log.Printf("failed to unmarshal: %v", rid)
			weque.SendError(w, 400, fmt.Sprintf("failed to unmarshal: %v", err))
			return
		}

		go (func() {
			events <- &Context{
				Webhook:         body,
				WebhookProvider: h.WebhookProvider(),
			}
			log.Printf("queued. delivery: %v, owner: %v, name: %v, pusher: %v", body.Delivery, body.Repository.Owner.Name, body.Repository.Name, body.Pusher.Name)
		})()

		w.Write([]byte(fmt.Sprintf("%v\n", body.Delivery)))
		log.Printf("respond. delivery: %v", body.Delivery)
	}
}
