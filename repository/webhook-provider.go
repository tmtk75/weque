package repository

type WebhookProvider interface {
	RepositoryURL(w *Webhook) string
	CommitURL(w *Webhook) string
	CompareURL(w *Webhook) string
	RefURL(w *Webhook) string
	PusherURL(w *Webhook) string
}
