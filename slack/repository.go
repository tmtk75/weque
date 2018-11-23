package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/repository"
)

const (
	KeySlackPayloadTemplateRepository = "notification.slack.payload_template_repository"
)

func NewIncomingWebhookRepository(w *repository.Webhook, u repository.WebhookProvider, exiterr error) (*IncomingWebhook, error) {
	templ := `
delivery:<{{ .RepositoryURL }}/settings/hooks|{{ .Delivery }}>
head_commit:<{{ .CommitURL }}?w=1|{{ .AfterShort }}>
ref:<{{ .RefURL }}|{{ .Repository.Owner.Name }}/{{ .Repository.Name }}:{{ .Ref }}>
compare:<{{ .CompareURL }}|{{ .BeforeShort }}...{{ .AfterShort }}>
pusher:<{{ .PusherURL }}|{{ .Pusher.Name }}>
elapsed: {{ .ElapsedTime }}
`

	if s := viper.GetString(KeySlackPayloadTemplateRepository); s != "" {
		templ = s
	}

	if len(w.Before) < 7 {
		return nil, errors.New("Before is less than 7")
	}
	if len(w.After) < 7 {
		return nil, errors.New("After is less than 7")
	}

	t := template.Must(template.New("").Parse(templ))
	text := bytes.NewBufferString("")
	err := t.Execute(text, struct {
		*repository.Webhook
		RepositoryURL string
		CommitURL     string
		CompareURL    string
		RefURL        string
		PusherURL     string
		AfterShort    string
		BeforeShort   string
		ElapsedTime   time.Duration
	}{
		Webhook:       w,
		RepositoryURL: u.RepositoryURL(w),
		CommitURL:     u.CommitURL(w),
		CompareURL:    u.CompareURL(w),
		RefURL:        u.RefURL(w),
		PusherURL:     u.PusherURL(w),
		BeforeShort:   w.Before[0:7],
		AfterShort:    w.After[0:7],
		ElapsedTime:   time.Duration(time.Now().Unix()-w.Repository.PushedAt) * time.Second,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed templating")
	}

	titletext := "ok"
	color := "good"
	if exiterr != nil {
		titletext = exiterr.Error()
		color = "danger"
	}

	wh := &IncomingWebhook{
		Channel:  viper.GetString(KeySlackChannelName),
		Username: fmt.Sprintf("webhook (%s)", w.Pusher.Name),
		Text:     titletext,
		Attachments: []Attachment{
			Attachment{
				AuthorName: u.Name(),
				AuthorIcon: u.IconURL(),
				Color:      color,
				Text:       text.String(),
			},
		},
	}

	return wh, nil
}

func PrintIncomingWebhookRepository(r *http.Request, path string, h repository.Handler, p repository.WebhookProvider, notify bool) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("%v", err)
	}

	wh, err := h.Unmarshal(r, b)
	if err != nil {
		log.Fatalf("%v", err)
	}
	iwh, err := NewIncomingWebhookRepository(wh, p, nil)

	s, err := json.Marshal(iwh)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", string(s))

	if notify {
		if err := Notify(iwh); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
