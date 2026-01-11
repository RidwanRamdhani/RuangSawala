package repositories

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) IsEmailTaken(email string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) CreateUser(username, email, hashedPassword string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO users (username, email) VALUES (?, ?)", username, email)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO auth (user_id, password_hash) VALUES (?, ?)", userID, hashedPassword)
	if err != nil {
		return err
	}

	return tx.Commit()
}
