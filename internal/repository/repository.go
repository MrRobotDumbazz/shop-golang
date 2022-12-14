package repository

import "database/sql"

type Repository struct {
	Client
	Auth
	Product
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Client:  newClientRepository(db),
		Auth:    newAuthRepository(db),
		Product: newProductRepository(db),
	}
}
