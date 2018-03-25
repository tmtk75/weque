package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/tmtk75/weque/repository"
)

type capture struct {
	request *http.Request
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

	err = Notify(&repository.Context{Webhook: &w, WebhookProvider: &repository.Github{}})

	assert.NoError(t, err)
	assert.NotNil(t, cap.request)

	p, _ := ioutil.ReadAll(cap.request.Body)
	var iwh IncomingWebhook
	err = json.Unmarshal(p, &iwh)
	assert.NoError(t, err)

	assert.Equal(t, "webhook (tmtk75)", iwh.Username)
	assert.Equal(t, "#api-test", iwh.Channel)
	assert.Equal(t, "something wrong", iwh.Text)
}
