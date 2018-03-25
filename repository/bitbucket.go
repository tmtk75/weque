package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Bitbucket struct {
}

const ENV_SECRET_TOKEN = "SECRET_TOKEN"

func (bb *Bitbucket) Verify(r *http.Request, body []byte) error {
	token := os.Getenv(ENV_SECRET_TOKEN)
	secret := r.URL.Query().Get("secret")
	if token != secret {
		return fmt.Errorf("the given secret token didn't match")
	}
	return nil
}

func (s *Bitbucket) IsPing(r *http.Request, body []byte) bool {
	return false
}

func (bb *Bitbucket) Unmarshal(r *http.Request, payload []byte) (*Webhook, error) {
	var body BitbucketWebhook
	if err := json.Unmarshal(payload, &body); err != nil {
		return nil, err
	}

	if _, ok := body.Push["changes"]; !ok {
		return nil, fmt.Errorf("push.changes is missing: %v", body)
	}
	if len(body.Push["changes"]) == 0 {
		return nil, fmt.Errorf("push.changes is empty: %v", body)
	}

	wb, err := body.Webhook()
	if err != nil {
		return nil, err
	}
	wb.Event = r.Header.Get("X-Event-Key")
	wb.Delivery = r.Header.Get("X-Request-UUID")
	return wb, nil
}

/*
func (bb *Bitbucket) NewKeyValue(r *http.Request, b []byte, wh *Webhook) (key string, val []byte, err error) {
	key = fmt.Sprintf("weque/bitbucket/%v/%v/webhooks/%v",
		wh.Repository.Owner.Name,
		wh.Repository.Name,
		wh.Delivery,
	)
	v, err := json.Marshal(struct {
		Webhook *Webhook          `json:"webhook"`
		Headers map[string]string `json:"headers"`
		Payload []byte            `json:"payload"`
	}{
		Webhook: wh,
		Headers: map[string]string{
			"X-Attempt-Number": r.Header.Get("X-Attempt-Number"),
			"X-Hook-UUID":      r.Header.Get("X-Hook-UUID"),
			"X-Event-Key":      r.Header.Get("X-Event-Key"),
			"X-Request-UUID":   r.Header.Get("X-Request-UUID"),
		},
		Payload: b,
	})
	if err != nil {
		return key, nil, err
	}
	val = v
	return
}

func (bb *Bitbucket) EventName(r *http.Request, b []byte, wh *Webhook) string {
	return fmt.Sprintf("bitbucket.%v", wh.Event)
}
*/
