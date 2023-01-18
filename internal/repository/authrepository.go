package repository

import (
	"database/sql"
	"errors"
	"shop/models"
)

type Auth interface {
	CreateSeller(*models.Seller) error
	GetUser(email, password string) (int, error)
	GetUserInID(id int) (*models.Seller, error)
}

type AuthRepository struct {
	db *sql.DB
}

var ErrRecordNotFound = errors.New("Record not found")

func newAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateSeller(s *models.Seller) error {
	if _, err := r.db.Exec("INSERT INTO shopdb.seller (email, password) VALUES (?, ?)", s.Email, s.Password); err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) GetUser(email, password string) (int, error) {
	s := &models.Seller{}
	err := r.db.QueryRow("SELECT id FROM shopdb.seller WHERE email = ?, password = ?",
		email, password).Scan(&s.ID)
	if err == sql.ErrNoRows {
		return 0, ErrRecordNotFound
	}
	if err != nil {
		return 0, err
	}
	return s.ID, nil
}

func (r *AuthRepository) GetUserInID(id int) (*models.Seller, error) {
	s := &models.Seller{}
	err := r.db.QueryRow("SELECT id FROM shopdb.seller WHERE id", id).Scan(&s.ID)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}
