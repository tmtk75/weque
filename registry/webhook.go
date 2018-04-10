package registry

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Event struct {
	ID      string  `json:"id"`
	Target  Target  `json:"target"`
	Request Request `json:"request"`
}

type Target struct {
	Repository string `json:"repository"`
	Digest     string `json:"digest"`
	URL        string `json:"url"`
	Tag        string `json:"tag"`
}

type Request struct {
	ID        string `json:"id"`
	Addr      string `json:"addr"`
	Host      string `json:"host"`
	UserAgent string `json:"useragent"`
}

type Webhook struct {
	Events []Event `json:"events"`
}

func Unmarshal(b []byte) (*Webhook, error) {
	var wh Webhook
	if err := json.Unmarshal(b, &wh); err != nil {
		return nil, err
	}
	return &wh, nil
}

func (e *Event) Env() []string {
	prefix := strings.ToUpper(viper.GetString("prefix"))
	f := func(a string, b interface{}) string { return fmt.Sprintf("%s%s=%v", prefix, a, b) }
	return []string{
		f("EVENT_ID", e.ID),
		f("REPOSITORY", e.Target.Repository),
		f("DIGEST", e.Target.Digest),
		f("URL", e.Target.URL),
		f("TAG", e.Target.Tag),
		f("REQUEST_ID", e.Request.ID),
		f("ADDR", e.Request.Addr),
		f("HOST", e.Request.Host),
		f("USER_AGENT", e.Request.UserAgent),
	}
}
