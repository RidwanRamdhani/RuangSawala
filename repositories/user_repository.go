package repositories

import (
	"database/sql"
	"time"

	"github.com/ruangsawala/backend/models"
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

func (r *UserRepository) GetUserByEmail(email string) (*models.User, *models.Auth, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at, a.password_hash, a.is_verified
		FROM users u
		JOIN auth a ON u.id = a.user_id
		WHERE u.email = ?
	`
	user := &models.User{}
	auth := &models.Auth{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
		&auth.PasswordHash, &auth.IsVerified,
	)
	if err != nil {
		return nil, nil, err
	}
	auth.UserID = user.ID
	return user, auth, nil
}

func (r *UserRepository) CreateUser(user *models.User, auth *models.Auth) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := tx.Exec("INSERT INTO users (username, email, created_at, updated_at) VALUES (?, ?, ?, ?)", user.Username, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(userID)
	auth.UserID = int(userID)

	_, err = tx.Exec("INSERT INTO auth (user_id, password_hash, is_verified) VALUES (?, ?, ?)", auth.UserID, auth.PasswordHash, auth.IsVerified)
	if err != nil {
		return err
	}

	return tx.Commit()
}
