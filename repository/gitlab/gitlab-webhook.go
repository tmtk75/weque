package gitlab

type GitlabWebhook struct {
	After   string `json:"after"`
	Before  string `json:"before"`
	Ref     string `json:"ref"`
	Project struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"project"`
	UserUsername string `json:"user_username"`
}
