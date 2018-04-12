package repositoryworker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/repository"
	"github.com/tmtk75/weque/repository/github"
	"github.com/tmtk75/weque/slack"
)

func TestRun(t *testing.T) {
	in := make(chan *repository.Context)
	out := Run(in)

	in <- &repository.Context{Webhook: &repository.Webhook{Delivery: "abcd"}}
	e := <-out
	assert.Equal(t, "abcd", e.Context.Webhook.Delivery)
	assert.Regexp(t, "failed to run", e.Err.Error())
}

func TestNotify(t *testing.T) {
	in := make(chan *Context)
	out := Notify(in)

	var called bool
	notifier = func(wh *slack.IncomingWebhook) error {
		called = true
		return nil
	}

	in <- &Context{Context: &repository.Context{
		Webhook:         &repository.Webhook{Before: "01234567", After: "01234567"},
		WebhookProvider: &github.Github{},
	}}
	err := <-out

	assert.True(t, called)
	assert.Nil(t, err)
}
