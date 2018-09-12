package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque/repository"
)

const (
	KeySlackGitlabIconURL = "notification.slack.gitlab_icon_url"
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
	wh.After = glwh.After
	wh.Before = glwh.Before
	wh.Ref = glwh.Ref
	wh.Repository.Name = glwh.Project.Name
	wh.Repository.Owner.Name = glwh.Project.Namespace
	wh.Pusher.Name = glwh.UserUsername
	return &wh, err
}

func (g *Gitlab) WebhookProvider() repository.WebhookProvider {
	return g
}

func (g *Gitlab) Name() string {
	return "gitlab"
}

func (g *Gitlab) IconURL() string {
	return viper.GetString(KeySlackGitlabIconURL)
}

func (g *Gitlab) RepositoryURL(w *repository.Webhook) string {
	return fmt.Sprintf("https://gitlab.com/%s/%s", w.Repository.Owner.Name, w.Repository.Name)
}

func (g *Gitlab) CommitURL(w *repository.Webhook) string {
	return fmt.Sprintf("%s/commits/%s", g.RepositoryURL(w), w.After)
}

func (g *Gitlab) CompareURL(w *repository.Webhook) string {
	return fmt.Sprintf("%s/compare/%s...%s", g.RepositoryURL(w), w.Before, w.After)
}

func (g *Gitlab) RefURL(w *repository.Webhook) string {
	s := strings.Split(w.Ref, "/")
	return fmt.Sprintf("%s/tree/%s", g.RepositoryURL(w), s[len(s)-1])
}

func (g *Gitlab) PusherURL(w *repository.Webhook) string {
	return fmt.Sprintf("https://gitlab.com/%s", w.Pusher.Name)
}
