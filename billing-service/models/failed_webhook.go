package models

type FailedWebhook struct {
	ID		 int    `db:"id" json:"id"`
	Event      string `db:"event" json:"event"`
	Payload    string `db:"payload" json:"payload"`
	RetryCount int    `db:"retry_count" json:"retry_count"`
}
