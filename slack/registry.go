package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/registry"
)

const (
	KeySlackPayloadTemplateRegistry = "notification.slack.payload_template_registry"
	KeySlackDockerIconURL           = "notification.slack.docker_icon_url"
)

func NewIncomingWebhookRegistry(e *registry.Event, exiterr error) (*IncomingWebhook, error) {
	templ := `
repository:{{ .Event.Target.Repository }}
tag:{{ .Event.Target.Tag }}
digest:{{ .Event.Target.Digest }}
url:{{ .Event.Target.URL }}
id:{{ .Event.ID }}
`
	if s := viper.GetString(KeySlackPayloadTemplateRegistry); s != "" {
		templ = s
	}

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
		Channel:  viper.GetString(KeySlackChannelName),
		Text:     titletext,
		Attachments: []Attachment{
			{
				AuthorName: e.Target.Repository,
				AuthorIcon: viper.GetString(KeySlackDockerIconURL),
				Color:      color,
				Text:       text.String(),
			},
		},
	}, nil
}

func PrintIncomingWebhookRegistry(path string, notify bool) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var wh registry.Webhook
	err = json.Unmarshal(b, &wh)
	if err != nil {
		log.Fatalf("%v", err)
	}
	iwh, err := NewIncomingWebhookRegistry(&wh.Events[0], nil)

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
