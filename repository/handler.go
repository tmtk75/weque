package repository

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tmtk75/weque"
)

type Handler interface {
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

	/*
	 * Returns key and value to be stored in a key-value store.
	 * Designed for consul KV.
	 */
	//NewKeyValue(r *http.Request, body []byte, wh *Webhook) (key string, val []byte, err error)

	/*
	 * Returns a name when firing event.
	 * Designed for consul event.
	 */
	//EventName(r *http.Request, body []byte, wh *Webhook) string
}

func NewHandler(wh Handler, events chan<- *Webhook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			weque.SendError(w, 500, fmt.Sprintf("%v", err))
			return
		}

		if err := wh.Verify(r, b); err != nil {
			log.Printf("%v", string(b))
			weque.SendError(w, 400, fmt.Sprintf("Failed to verify: %v", err))
			return
		}

		if wh.IsPing(r, b) {
			msg := fmt.Sprintf("Received pings: %v", r.RequestURI)
			log.Print(msg)
			w.WriteHeader(200)
			w.Write([]byte(msg))
			return
		}

		body, err := wh.Unmarshal(r, b)
		if err != nil {
			weque.SendError(w, 400, fmt.Sprintf("%v", err))
			return
		}

		go (func() {
			events <- body
			log.Printf("queued. delivery: %v, owner: %v, name: %v, pusher: %v", body.Delivery, body.Repository.Owner.Name, body.Repository.Name, body.Pusher.Name)
		})()

		w.Write([]byte(fmt.Sprintf("%v\n", body.Delivery)))
		log.Printf("respond %v", body.Delivery)
	}
}
