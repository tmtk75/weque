package slack

import (
	"log"

	"github.com/tmtk75/weque/registry"
)

func NotifyRegistry(e *registry.Event, err error) error {
	wh, err := newRegistryIncomingWebhook(e, err)
	if err != nil {
		return err
	}
	if err := request(wh); err != nil {
		return err
	}
	log.Printf("notified to slack channel, %v", wh.Channel)
	return nil
}

func newRegistryIncomingWebhook(e *registry.Event, err error) (*IncomingWebhook, error) {
	var (
		color = "good"
		text  = "ok"
	)
	if err != nil {
		color = "danger"
		text = err.Error()
	}
	return &IncomingWebhook{
		Username: "webhook",
		Channel:  "#api-test",
		Text:     text,
		Attachments: []Attachment{
			{
				AuthorName: e.Target.Repository,
				AuthorIcon: "http://calvintrobinson.s3.amazonaws.com/wp-content/uploads/harbor-logo2.png",
				Color:      color,
				Text:       e.Target.Repository + ":" + e.Target.Tag,
			},
		},
	}, nil
}
