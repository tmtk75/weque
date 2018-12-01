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
		e("ACTION", wh.Action),
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
		e("RELEASE_TAG_NAME", wh.Release.TagName),
	}
}
