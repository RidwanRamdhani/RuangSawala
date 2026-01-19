package repositories

import (
	"database/sql"
	"errors"

	"github.com/ruangsawala/backend/models"
)

type UserInterestRepository struct {
	db *sql.DB
}

func NewUserInterestRepository(db *sql.DB) *UserInterestRepository {
	return &UserInterestRepository{db: db}
}

func (r *UserInterestRepository) AddUserInterest(userID int, interestID int) error {
	if userID <= 0 {
		return errors.New("invalid userID")
	}
	if interestID <= 0 {
		return errors.New("invalid interestID")
	}
	query := `INSERT INTO user_interests (user_id, interest_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, userID, interestID)
	return err
}

func (r *UserInterestRepository) RemoveUserInterest(userID int, interestID int) error {
	if userID <= 0 {
		return errors.New("invalid userID")
	}
	if interestID <= 0 {
		return errors.New("invalid interestID")
	}
	query := `DELETE FROM user_interests WHERE user_id = ? AND interest_id = ?`
	_, err := r.db.Exec(query, userID, interestID)
	return err
}

func (r *UserInterestRepository) GetUserInterests(userID int) ([]models.Interest, error) {
	if userID <= 0 {
		return nil, errors.New("invalid userID")
	}
	query := `SELECT i.id, i.name FROM interests i JOIN user_interests ui ON i.id = ui.interest_id WHERE ui.user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interests []models.Interest
	for rows.Next() {
		var interest models.Interest
		err := rows.Scan(&interest.ID, &interest.Name)
		if err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return interests, nil
}

func (r *UserInterestRepository) GetAllInterests() ([]models.Interest, error) {
	query := `SELECT id, name FROM interests`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interests []models.Interest
	for rows.Next() {
		var interest models.Interest
		err := rows.Scan(&interest.ID, &interest.Name)
		if err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return interests, nil
}
