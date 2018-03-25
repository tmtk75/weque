package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tmtk75/weque/repository"
)

const (
	KEY_SLACK_URL    = "notification.slack.incoming_webhook_url"
	KEY_CHANNEL_NAME = "notification.slack.channel_name"
)

func init() {
	viper.BindEnv(KEY_SLACK_URL, "SLACK_URL")
	viper.BindEnv(KEY_CHANNEL_NAME, "SLACK_CHANNEL_NAME")
}

func Notify(c *repository.Context) error {
	wh, err := newIncomingWebhook(c.Webhook, c.WebhookProvider)
	if err != nil {
		return err
	}
	log.Printf("notified to slack channel, %v", wh.Channel)
	return request(wh)
}

var newClient = func(c *http.Client) *http.Client {
	return c
}

func request(wh *IncomingWebhook) error {
	b, err := json.Marshal(wh)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %v", wh)
	}

	r := bytes.NewBuffer(b)
	url := viper.GetString(KEY_SLACK_URL)
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to build a request for slack notification")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	req = req.WithContext(ctx)
	defer cancel()

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second, // Timeout to connect
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
	client = newClient(client)

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to send request for slack notification")
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return errors.Errorf("received non 2xx status. err: %v", res.Status)
	}

	return nil
}
