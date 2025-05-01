package models

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password,omitempty"`
	Role     string `db:"role" json:"role"`
}
