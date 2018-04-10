package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/repository"
	bb "github.com/tmtk75/weque/repository/bitbucket"
	gh "github.com/tmtk75/weque/repository/github"
)

type capture struct {
	request *http.Request
}

type incomingWebHook struct {
	err                                           error
	username, channel, text                       string
	authorIcon, authorName, color, attachmentText string
}

func (c *capture) RoundTrip(r *http.Request) (*http.Response, error) {
	c.request = r
	return &http.Response{
		StatusCode: 201,
		Body:       ioutil.NopCloser(nil),
	}, nil
}

func TestMain(m *testing.M) {
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	c := m.Run()
	os.Exit(c)
}

func TestNotify(t *testing.T) {
	b, err := ioutil.ReadFile("../github/payload.json")
	if err != nil {
		panic(err)
	}

	var w repository.Webhook
	err = json.Unmarshal(b, &w)
	if err != nil {
		panic(err)
	}
	w.Delivery = "delivery-test-foobar"

	//
	cap := &capture{}
	newClient = func(c *http.Client) *http.Client {
		c.Transport = cap
		return c
	}

	run := func(expects []incomingWebHook, provider repository.WebhookProvider) {
		for _, e := range expects {
			ctx := &repository.Context{Webhook: &w, WebhookProvider: provider}
			err = Notify(ctx, e.err)

			assert.NoError(t, err)
			assert.NotNil(t, cap.request)

			p, _ := ioutil.ReadAll(cap.request.Body)
			var iwh IncomingWebhook
			err = json.Unmarshal(p, &iwh)
			assert.NoError(t, err)

			assert.Equal(t, e.username, iwh.Username)
			assert.Equal(t, e.channel, iwh.Channel)
			assert.Equal(t, e.text, iwh.Text)
			assert.Len(t, iwh.Attachments, 1)
			assert.Equal(t, e.authorIcon, iwh.Attachments[0].AuthorIcon)
			assert.Equal(t, e.authorName, iwh.Attachments[0].AuthorName)
			assert.Equal(t, e.color, iwh.Attachments[0].Color)
			assert.Regexp(t, e.attachmentText, iwh.Attachments[0].Text)
		}
	}

	// GitHub
	t.Run("github", func(t *testing.T) {
		icon := "http://cdn.flaticon.com/png/256/25231.png"
		expects := []incomingWebHook{
			{
				err:      nil,
				username: "webhook (tmtk75)", channel: "#api-test", text: "ok",
				authorIcon: icon, authorName: "github", color: "good", attachmentText: "\ndelivery:<",
			},
			{
				err:      errors.Errorf("failed to run for github"),
				username: "webhook (tmtk75)", channel: "#api-test", text: "failed to run for github",
				authorIcon: icon, authorName: "github", color: "danger", attachmentText: "\ndelivery:<",
			},
		}
		run(expects, &gh.Github{})
	})

	// Bitbucket
	t.Run("bitbucket", func(t *testing.T) {
		icon := "https://www.atlassian.com/dam/jcr:e2a6f06f-b3d5-4002-aed3-73539c56a2eb/Bitbucket@2x-blue.png"
		expects := []incomingWebHook{
			{
				err:      nil,
				username: "webhook (tmtk75)", channel: "#api-test", text: "ok",
				authorIcon: icon, authorName: "bitbucket", color: "good", attachmentText: "\ndelivery:<",
			},
			{
				err:      errors.Errorf("failed to run for bitbucket"),
				username: "webhook (tmtk75)", channel: "#api-test", text: "failed to run for bitbucket",
				authorIcon: icon, authorName: "bitbucket", color: "danger", attachmentText: "\ndelivery:<",
			},
		}
		run(expects, &bb.Bitbucket{})
	})

}
