package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque/repository"
)

const (
	KeySlackBitbucketIconURL = "notification.slack.bitbucket_icon_url"
)

type Bitbucket struct {
}

func (bb *Bitbucket) RequestID(r *http.Request) string {
	return r.Header.Get("X-Request-UUID")
}

func (bb *Bitbucket) Verify(r *http.Request, body []byte) error {
	token := viper.GetString(repository.KeySecretToken)
	secret := r.URL.Query().Get("secret")
	if token != secret {
		return fmt.Errorf("the given secret token didn't match")
	}
	return nil
}

func (s *Bitbucket) IsPing(r *http.Request, body []byte) bool {
	return false
}

func (bb *Bitbucket) Unmarshal(r *http.Request, payload []byte) (*repository.Webhook, error) {
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
func (bb *Bitbucket) NewKeyValue(r *http.Request, b []byte, wh *repository.Webhook) (key string, val []byte, err error) {
	key = fmt.Sprintf("weque/bitbucket/%v/%v/webhooks/%v",
		wh.Repository.Owner.Name,
		wh.Repository.Name,
		wh.Delivery,
	)
	v, err := json.Marshal(struct {
		repository.Webhook *repository.Webhook          `json:"webhook"`
		Headers map[string]string `json:"headers"`
		Payload []byte            `json:"payload"`
	}{
		repository.Webhook: wh,
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

func (bb *Bitbucket) EventName(r *http.Request, b []byte, wh *repository.Webhook) string {
	return fmt.Sprintf("bitbucket.%v", wh.Event)
}
*/

func (s *Bitbucket) WebhookProvider() repository.WebhookProvider {
	return s
}

func (s *Bitbucket) Name() string {
	return "bitbucket"
}

func (s *Bitbucket) IconURL() string {
	return viper.GetString(KeySlackBitbucketIconURL)
}

func (s *Bitbucket) RepositoryURL(w *repository.Webhook) string {
	return fmt.Sprintf("https://bitbucket.org/%s/%s", w.Repository.Owner.Name, w.Repository.Name)
}

func (s *Bitbucket) CommitURL(w *repository.Webhook) string {
	return fmt.Sprintf("%s/commits/%s", s.RepositoryURL(w), w.After)
}

func (s *Bitbucket) CompareURL(w *repository.Webhook) string {
	return fmt.Sprintf("%s/branches/compare/%s..%s", s.RepositoryURL(w), w.Before, w.After)
}

func (s *Bitbucket) RefURL(w *repository.Webhook) string {
	ref := strings.TrimPrefix(w.Ref, "refs/heads/")
	return fmt.Sprintf("%s/branch/%s", s.RepositoryURL(w), ref)
}

func (s *Bitbucket) PusherURL(w *repository.Webhook) string {
	return fmt.Sprintf("https://bitbucket.org/%s", w.Pusher.Name)
}
