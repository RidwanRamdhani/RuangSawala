package models

type Auth struct {
	UserID       int    `json:"user_id" db:"user_id"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	IsVerified   bool   `json:"is_verified" db:"is_verified"`
}
