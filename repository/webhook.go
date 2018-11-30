package repository

/*
 * Normalized webhook payload.
 *
 * This structure is for GitHub webhook.
 * https://developer.github.com/webhooks/
 *
 * This doesn't contain all fields of GitHub webhook payload
 * because it's designed to propagate event into other KVS, queue, pub/sub, etc.
 * which event payload max size is not so much, for example, consul is 512 bytes.
 */
type Webhook struct {
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
		PushedAt interface{} `json:"pushed_at"` // 2018-11-23T10:10:08Z (ping, release) or 1543592283 (pushed)
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
