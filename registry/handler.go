package registry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	KeyHandlerScript = "handlers.registry"
)

func RegistryHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read: %v", err)
		return
	}

	var body RegistryWebhookBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Printf("failed to unmarshal: %v", err)
		return
	}

	log.Printf("%v", r)
	w.Write([]byte("ok"))
}

type RegistryWebhookBody struct {
}
