package registry

import "encoding/json"

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