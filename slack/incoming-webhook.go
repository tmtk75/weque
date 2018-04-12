package slack

// https://kii.slack.com/services/B0JGN99RN?added=1
type IncomingWebhook struct {
	Channel     string       `json:"channel"` // e.g) #api-test
	Username    string       `json:"username"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	AuthorName string `json:"author_name"` // e.g) github
	AuthorIcon string `json:"author_icon"` // e.g) "http://cdn.flaticon.com/png/256/25231.png"
	Color      string `json:"color"`       // e.g) good, warn, danger, #00ff00
	Text       string `json:"text"`
}
