package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/registry"
)

func TestNotifyRegistry(t *testing.T) {
	b, err := ioutil.ReadFile("../registry/payload.json")
	assert.NoError(t, err)

	wh, err := registry.Unmarshal(b)
	assert.NoError(t, err)
	event := &wh.Events[0]

	//
	cap := &capture{}
	newClient = func(c *http.Client) *http.Client {
		c.Transport = cap
		return c
	}

	//
	icon := "http://calvintrobinson.s3.amazonaws.com/wp-content/uploads/harbor-logo2.png"
	expects := []incomingWebHook{
		{
			err:      nil,
			username: "webhook", channel: "#api-test", text: "ok",
			authorIcon: icon, authorName: "alpine", color: "good",
			attachmentText: "repository:alpine\ntag:3.6",
		},
		{
			err:      errors.Errorf("failed to run for github"),
			username: "webhook", channel: "#api-test", text: "failed to run for github",
			authorIcon: icon, authorName: "alpine", color: "danger",
			attachmentText: "repository:alpine\ntag:3.6",
		},
	}

	for _, e := range expects {
		wh, _ := NewRegistryIncomingWebhook(event, e.err)
		Notify(wh)
		defer cap.request.Body.Close()
		body, _ := ioutil.ReadAll(cap.request.Body)
		//t.Logf("%v", string(body))

		var iwh IncomingWebhook
		err := json.Unmarshal(body, &iwh)
		assert.NoError(t, err)
		//t.Logf("%v", iwh)

		assert.Equal(t, e.username, iwh.Username, "username")
		assert.Equal(t, e.channel, iwh.Channel)
		assert.Equal(t, e.text, iwh.Text, "text")
		assert.Len(t, iwh.Attachments, 1)
		assert.Equal(t, e.authorIcon, iwh.Attachments[0].AuthorIcon, "icon")
		assert.Equal(t, e.authorName, iwh.Attachments[0].AuthorName)
		assert.Equal(t, e.color, iwh.Attachments[0].Color)
		assert.Regexp(t, e.attachmentText, iwh.Attachments[0].Text)
	}
}
