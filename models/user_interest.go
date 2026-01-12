package models

type UserInterest struct {
	UserID     int `json:"user_id" db:"user_id"`
	InterestID int `json:"interest_id" db:"interest_id"`
}
