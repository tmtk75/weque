package slack

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/repository"
)

// https://kii.slack.com/services/B0JGN99RN?added=1
type IncomingWebhook struct {
	Channel     string       `json:"channel"` // e.g) #api-test
	Username    string       `json:"username"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	AuthorName string `json:"author_name"` // e.g) github
	AuthorIcon string `json:"author_icon"` // e.g) "http://cdn.flaticon.com/png/256/25231.png"
	Color      string `json:"color"`       // e.g) good, warn, danger, #00ff00
	Text       string `json:"text"`
}

func NewIncomingWebhook(w *repository.Webhook, u repository.WebhookProvider, exiterr error) (*IncomingWebhook, error) {
	templ := `
delivery:<{{ .RepositoryURL }}/settings/hooks|{{ .Delivery }}>
head_commit:<{{ .CommitURL }}?w=1|{{ .AfterShort }}>
ref:<{{ .RefURL }}|{{ .Repository.Owner.Name }}/{{ .Repository.Name }}:{{ .Ref }}>
compare:<{{ .CompareURL }}|{{ .BeforeShort }}...{{ .AfterShort }}>
pusher:<{{ .PusherURL }}|{{ .Pusher.Name }}>
status:$status
`

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
	}{
		Webhook:       w,
		RepositoryURL: u.RepositoryURL(w),
		CommitURL:     u.CommitURL(w),
		CompareURL:    u.CompareURL(w),
		RefURL:        u.RefURL(w),
		PusherURL:     u.PusherURL(w),
		BeforeShort:   w.Before[0:7],
		AfterShort:    w.After[0:7],
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
		Channel:  viper.GetString(KEY_CHANNEL_NAME),
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
