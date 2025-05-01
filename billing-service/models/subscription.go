package models

import "time"

type Subscription struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Plan      string    `db:"plan" json:"plan"`
	Status    string    `db:"status" json:"status"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
}
