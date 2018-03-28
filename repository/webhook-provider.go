package repository

type WebhookProvider interface {
	Name() string
	IconURL() string
	RepositoryURL(w *Webhook) string
	CommitURL(w *Webhook) string
	CompareURL(w *Webhook) string
	RefURL(w *Webhook) string
	PusherURL(w *Webhook) string
}
