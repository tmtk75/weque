package repository

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func (wh *Webhook) Env() []string {
	prefix := strings.ToUpper(viper.GetString("prefix"))
	e := func(a string, b interface{}) string { return fmt.Sprintf("%s%s=%v", prefix, a, b) }
	return []string{
		e("REPOSITORY_NAME", wh.Repository.Name),
		e("OWNER_NAME", wh.Repository.Owner.Name),
		e("EVENT", wh.Event),
		e("DELIVERY", wh.Delivery),
		e("REF", wh.Ref),
		e("AFTER", wh.After),
		e("BEFORE", wh.Before),
		e("CREATED", wh.Created),
		e("DELETED", wh.Deleted),
		e("PUSHER_NAME", wh.Pusher.Name),
	}
}

/*
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"repository"`
	Event    string `json:"event"`
	Delivery string `json:"delivery"`
	Ref      string `json:"ref"`
	After    string `json:"after"`
	Before   string `json:"before"`
	Created  bool   `json:"created"`
	Deleted  bool   `json:"deleted"`
	//Head_commit map[string]interface{} `json:"head_commit,omitempty"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher,omitempty"`
*/
