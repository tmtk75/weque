package bitbucket

import (
	"fmt"

	"github.com/tmtk75/weque/repository"
)

/*
 * This structure is for BitBucket webhook.
 * https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html
 */
type BitbucketWebhook struct {
	Actor struct {
		Username string `json:"username"`
	} `json:"actor"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Type     string `json:type`
			Username string `json:username`
		} `json:"owner"`
	} `json:"repository"`
	Push map[string][]struct { // changes
		New struct {
			Name   string `json:"name"`
			Type   string `json:"type"`
			Target struct {
				Hash string `json:"hash"`
			} `json:"target"`
		} `json:"new"`
		Old struct {
			Target struct {
				Hash string `json:"hash"`
			} `json:"target"`
		} `json:"old"`
		Created   bool `json:"created"`
		Truncated bool `json:"truncated"`
	} `json:"push"`
}

/*
 * Returns a repository.Webhook converting itself.
 */
func (body *BitbucketWebhook) Webhook() (*repository.Webhook, error) {
	var wb repository.Webhook
	ch := body.Push["changes"][0]
	switch ch.New.Type {
	case "branch":
		wb.Ref = "refs/heads/" + ch.New.Name // Normalize
	case "tag":
		wb.Ref = "refs/tags/" + ch.New.Name // Normalize
	default:
		return nil, fmt.Errorf("unknown type: %v", ch.New.Type)
	}
	wb.Repository.Name = body.Repository.Name
	wb.Repository.Owner.Name = body.Repository.Owner.Username
	wb.Pusher.Name = body.Actor.Username
	wb.Created = ch.Created
	wb.Deleted = ch.Truncated
	wb.After = ch.New.Target.Hash
	wb.Before = ch.Old.Target.Hash
	return &wb, nil
}
