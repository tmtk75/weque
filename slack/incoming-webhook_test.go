package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/repository"
)

func TestNewIncomingWebhook(t *testing.T) {
	w := &repository.Webhook{Before: "0123456", After: "0123456"}
	g := &repository.Github{}
	iw, err := newIncomingWebhook(w, g, nil)

	assert.NoError(t, err)
	expected := "\ndelivery:<https://github.com///settings/hooks|>\nhead_commit:<https://github.com///commit/0123456?w=1|0123456>\nref:<https://github.com///tree/|/:>\ncompare:<https://github.com///compare/0123456...0123456?w=1|0123456...0123456>\npusher:<https://github.com/|>\nstatus:$status\n"
	assert.Equal(t, expected, iw.Attachments[0].Text)
}
