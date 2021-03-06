package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/tmtk75/weque/repository"
)

const (
	KeySlackGithubIconURL = "notification.slack.github_icon_url"
)

type Github struct {
}

func (s *Github) RequestID(r *http.Request) string {
	return r.Header.Get("X-Github-Delivery")
}

func (s *Github) Verify(r *http.Request, body []byte) error {
	sign := r.Header.Get("X-Hub-Signature")
	return Verify(sign, body)
}

func Verify(sign string, b []byte) error {
	log.Printf("sign: %v, len: %v", sign, len(b))
	if !strings.HasPrefix(sign, "sha1=") {
		s := fmt.Sprintf("unknown hash algorithm: %v", sign)
		log.Println(s)
		return errors.New(s)
	}

	expected, _ := hex.DecodeString(string(sign[4+1 : len(sign)])) // 4+1 is to skip `sha1=`
	expected = []byte(expected)

	token := viper.GetString(repository.KeySecretToken)
	mac := hmac.New(sha1.New, []byte(token))
	mac.Write(b)
	actual := mac.Sum(nil)

	if !hmac.Equal(actual, expected) {
		return errors.New("not match")
	}

	return nil
}

func (s *Github) IsPing(r *http.Request, body []byte) bool {
	e := r.Header.Get("X-Github-Event")
	return e == "ping"
}

func (s *Github) Unmarshal(r *http.Request, b []byte) (*repository.Webhook, error) {
	ctype := r.Header.Get("content-type")
	switch strings.ToLower(ctype) {
	case "application/x-www-form-urlencoded":
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		p := r.Form["payload"]
		if len(p) != 1 {
			return nil, fmt.Errorf("unexpected payload for %v: %v", ctype, p)
		}
		log.Printf("[debug] github payload: %v", repository.Shorten(p[0], 256))
		b = []byte(r.Form["payload"][0])
	case "application/json":
		// NOP
	default:
		return nil, fmt.Errorf("unsupported content-type: %v", ctype)
	}

	var body repository.Webhook
	if err := json.Unmarshal(b, &body); err != nil {
		return nil, err
	}
	body.Event = r.Header.Get("X-Github-Event")
	body.Delivery = r.Header.Get("X-Github-Delivery")
	return &body, nil
}

func (s *Github) WebhookProvider() repository.WebhookProvider {
	return s
}

func (s *Github) Name() string {
	return "github"
}

func (s *Github) IconURL() string {
	return viper.GetString(KeySlackGithubIconURL)
}

func (s *Github) RepositoryURL(w *repository.Webhook) string {
	return fmt.Sprintf("https://github.com/%s/%s", w.Repository.Owner.Name, w.Repository.Name)
}

func (s *Github) CommitURL(w *repository.Webhook) string {
	return fmt.Sprintf("%s/commit/%s", s.RepositoryURL(w), w.After)
}

func (s *Github) CompareURL(w *repository.Webhook) string {
	return fmt.Sprintf("%s/compare/%s...%s?w=1", s.RepositoryURL(w), w.Before, w.After)
}

func (s *Github) RefURL(w *repository.Webhook) string {
	return fmt.Sprintf("%s/tree/%s", s.RepositoryURL(w), w.Ref)
}

func (s *Github) PusherURL(w *repository.Webhook) string {
	return fmt.Sprintf("https://github.com/%s", w.Pusher.Name)
}
