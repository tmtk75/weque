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

/*
{
   "events": [
      {
         "id": "320678d8-ca14-430f-8bb6-4ca139cd83f7",
         "timestamp": "2016-03-09T14:44:26.402973972-08:00",
         "action": "pull",
         "target": {
            "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
            "size": 708,
            "digest": "sha256:fea8895f450959fa676bcc1df0611ea93823a735a01205fd8622846041d0c7cf",
            "length": 708,
            "repository": "hello-world",
            "url": "http://192.168.100.227:5000/v2/hello-world/manifests/sha256:fea8895f450959fa676bcc1df0611ea93823a735a01205fd8622846041d0c7cf",
            "tag": "latest"
         },
         "request": {
            "id": "6df24a34-0959-4923-81ca-14f09767db19",
            "addr": "192.168.64.11:42961",
            "host": "192.168.100.227:5000",
            "method": "GET",
            "useragent": "curl/7.38.0"
         },
         "actor": {
            "name": "docker-hub"
	 },
         "source": {
            "addr": "xtal.local:5000",
            "instanceID": "a53db899-3b4b-4a62-a067-8dd013beaca4"
         }
      }
   ]
}
*/
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
