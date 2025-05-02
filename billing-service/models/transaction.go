package models

type Transaction struct {
	TransactionID   string `db:"transaction_id" json:"transaction_id"` // bug: used `id` instaed of `transaction_id`
	TxRef           string `db:"tx_ref" json:"tx_ref"`
	Amount          string `db:"amount" json:"amount"`
	Currency        string `db:"currency" json:"currency"`
	Status          string `db:"status" json:"status"`
	CustomerEmail   string `db:"customer_email" json:"customer_email"`
	TransactionDate string `db:"transaction_date" json:"transaction_date"`
	UserID          string `db:"user_id" json:"user_id"`
}
