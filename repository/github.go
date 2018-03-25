package repository

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
)

type Github struct {
}

func (s *Github) Verify(r *http.Request, body []byte) error {
	sign := r.Header.Get("X-Hub-Signature")
	return Verify(sign, body)
}

func Verify(sign string, b []byte) error {
	if !strings.HasPrefix(sign, "sha1=") {
		s := fmt.Sprintf("unknown hash algorithm: %v", sign)
		log.Println(s)
		return errors.New(s)
	}

	expected, _ := hex.DecodeString(string(sign[4+1 : len(sign)])) // 4+1 is to skip `sha1=`
	expected = []byte(expected)

	token := viper.GetString(KeySecretToken)
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

func (s *Github) Unmarshal(r *http.Request, b []byte) (*Webhook, error) {
	ctype := r.Header.Get("content-type")
	switch strings.ToLower(ctype) {
	case "application/x-www-form-urlencoded":
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		p := r.Form["payload"]
		if len(p) != 1 {
			return nil, fmt.Errorf("unexpected payload: %v", p)
		}
		b = []byte(r.Form["payload"][0])
	case "application/json":
		// NOP
	default:
		return nil, fmt.Errorf("unsupported content-type: %v", ctype)
	}

	var body Webhook
	if err := json.Unmarshal(b, &body); err != nil {
		return nil, err
	}
	body.Event = r.Header.Get("X-Github-Event")
	body.Delivery = r.Header.Get("X-Github-Delivery")
	return &body, nil
}

/* suspended for consul
func (s *Github) NewKeyValue(r *http.Request, b []byte, wh *Webhook) (key string, val []byte, err error) {
	key = fmt.Sprintf("weque/github/%v/%v/webhooks/%v",
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
		Payload: b,
		Headers: map[string]string{
			"X-Github-Event":   r.Header.Get("X-Github-Event"),
			"X-Github-Deliver": r.Header.Get("X-Github-Delivery"),
		},
	})
	if err != nil {
		return key, nil, err
	}
	val = v
	return
}

func (s *Github) EventName(r *http.Request, b []byte, wh *Webhook) string {
	return fmt.Sprintf("github.%v", wh.Event)
}
*/

func (s *Github) WebhookProvider() WebhookProvider {
	return s
}

func (s *Github) RepositoryURL(w *Webhook) string {
	return fmt.Sprintf("https://github.com/%s/%s", w.Repository.Owner.Name, w.Repository.Name)
}

func (s *Github) CommitURL(w *Webhook) string {
	return fmt.Sprintf("%s/commit/%s", s.RepositoryURL(w), w.After)
}

func (s *Github) CompareURL(w *Webhook) string {
	return fmt.Sprintf("%s/compare/%s...%s?w=1", s.RepositoryURL(w), w.Before, w.After)
}

func (s *Github) RefURL(w *Webhook) string {
	return fmt.Sprintf("%s/tree/%s", s.RepositoryURL(w), w.Ref)
}

func (s *Github) PusherURL(w *Webhook) string {
	return fmt.Sprintf("https://github.com/%s", w.Pusher.Name)
}
