package registryworker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/registry"
	"github.com/tmtk75/weque/slack"
)

func TestRun(t *testing.T) {
	in := make(chan *registry.Webhook)
	out := Run(in)

	in <- &registry.Webhook{Events: []registry.Event{{ID: "abcd"}}}
	e := <-out
	assert.Equal(t, "abcd", e.ID)
	assert.Regexp(t, "failed to run", e.Err.Error())
}

func TestNotify(t *testing.T) {
	in := make(chan *Event)
	out := Notify(in)

	var called bool
	notifier = func(wh *slack.IncomingWebhook) error {
		called = true
		return nil
	}

	in <- &Event{Event: &registry.Event{
		Target: registry.Target{
			Repository: "",
		},
	}}
	err := <-out

	assert.True(t, called)
	assert.Nil(t, err)
}
