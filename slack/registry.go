package slack

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
	"github.com/tmtk75/weque/registry"
)

func NewIncomingWebhookRegistry(e *registry.Event, exiterr error) (*IncomingWebhook, error) {
	templ := `
repository:{{ .Event.Target.Repository }}
tag:{{ .Event.Target.Tag }}
digest:{{ .Event.Target.Digest }}
url:{{ .Event.Target.URL }}
id:{{ .Event.ID }}
`

	t := template.Must(template.New("").Parse(templ))
	text := bytes.NewBufferString("")
	err := t.Execute(text, struct {
		*registry.Event
	}{
		Event: e,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed templating")
	}

	var (
		color     = "good"
		titletext = "ok"
	)
	if exiterr != nil {
		color = "danger"
		titletext = exiterr.Error()
	}
	return &IncomingWebhook{
		Username: "webhook",
		Channel:  "#api-test",
		Text:     titletext,
		Attachments: []Attachment{
			{
				AuthorName: e.Target.Repository,
				AuthorIcon: "http://calvintrobinson.s3.amazonaws.com/wp-content/uploads/harbor-logo2.png",
				Color:      color,
				Text:       text.String(),
			},
		},
	}, nil
}
