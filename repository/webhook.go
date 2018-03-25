package repository

import (
	"encoding/json"
	"fmt"
)

/*
 * This structure is for GitHub webhook.
 * https://developer.github.com/webhooks/
 *
 * This doesn't contain all fields of GitHub webhook payload
 * because it's designed to propagate event in consul
 * which event payload size is 512 bytes.
 */
type Webhook struct {
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
}

/*
 * Returns bytes in JSON.
 */
func (b *Webhook) Bytes() []byte {
	r, _ := json.Marshal(b)
	return r
}

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
 * Returns a Webhook converting itself.
 */
func (body *BitbucketWebhook) Webhook() (*Webhook, error) {
	var wb Webhook
	ch := body.Push["changes"][0]
	switch ch.New.Type {
	case "branch":
		wb.Ref = "refs/heads/" + ch.New.Name
	case "tag":
		wb.Ref = "refs/tags/" + ch.New.Name
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
