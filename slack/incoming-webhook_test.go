package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/repository"
	gh "github.com/tmtk75/weque/repository/github"
)

func TestNewIncomingWebhook(t *testing.T) {
	w := &repository.Webhook{Before: "0123456", After: "0123456"}
	g := &gh.Github{}
	iw, err := NewIncomingWebhook(w, g, nil)

	assert.NoError(t, err)
	expected := "\ndelivery:<https://github.com///settings/hooks|>\nhead_commit:<https://github.com///commit/0123456?w=1|0123456>\nref:<https://github.com///tree/|/:>\ncompare:<https://github.com///compare/0123456...0123456?w=1|0123456...0123456>\npusher:<https://github.com/|>\nstatus:$status\n"
	assert.Equal(t, expected, iw.Attachments[0].Text)
}

func TestNewIncomingWebhookBeforeAfter(t *testing.T) {
	w := &repository.Webhook{Before: "0123456"}
	_, err := NewIncomingWebhook(w, &gh.Github{}, nil)
	assert.Error(t, err)

	w = &repository.Webhook{After: "0123456"}
	_, err = NewIncomingWebhook(w, &gh.Github{}, nil)
	assert.Error(t, err)
}
