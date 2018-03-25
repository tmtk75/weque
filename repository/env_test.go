package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookEnv(t *testing.T) {
	wh := &Webhook{}
	expected := []string{
		"REPOSITORY_NAME=",
		"OWNER_NAME=",
		"EVENT=",
		"DELIVERY=",
		"REF=",
		"AFTER=",
		"BEFORE=",
		"CREATED=false",
		"DELETED=false",
		"PUSHER_NAME=",
	}

	e := wh.Env()
	for i, j := range expected {
		assert.Equal(t, e[i], j)
	}
}
