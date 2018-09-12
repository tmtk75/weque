package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque/repository"
)

type Gitlab struct {
}

func (g *Gitlab) RequestID(r *http.Request) string {
	return "id is not given for gitlab"
}

func (g *Gitlab) Verify(r *http.Request, body []byte) error {
	token := viper.GetString(repository.KeySecretToken)
	secret := r.Header.Get("X-Gitlab-Token")
	if token != secret {
		return fmt.Errorf("the given secret token didn't match")
	}
	return nil
}

func (g *Gitlab) IsPing(r *http.Request, body []byte) bool {
	e := r.Header.Get("X-Gitlab-Event")
	return strings.ToLower(e) == "push hook"
}

func (g *Gitlab) Unmarshal(r *http.Request, body []byte) (*repository.Webhook, error) {
	var glwh GitlabWebhook
	err := json.Unmarshal(body, &glwh)

	var wh repository.Webhook
	return &wh, err
}

func (g *Gitlab) WebhookProvider() repository.WebhookProvider {
	return g
}

func (g *Gitlab) Name() string {
	return ""
}

func (g *Gitlab) IconURL() string {
	return ""
}

func (g *Gitlab) RepositoryURL(w *repository.Webhook) string {
	return ""
}

func (g *Gitlab) CommitURL(w *repository.Webhook) string {
	return ""
}

func (g *Gitlab) CompareURL(w *repository.Webhook) string {
	return ""
}

func (g *Gitlab) RefURL(w *repository.Webhook) string {
	return ""
}

func (g *Gitlab) PusherURL(w *repository.Webhook) string {
	return ""
}
