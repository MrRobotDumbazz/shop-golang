package repository

import (
	"database/sql"
)

type Repository struct {
	Seller
	Auth
	Product
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Seller:  newSellerRepository(db),
		Auth:    newAuthRepository(db),
		Product: newProductRepository(db),
	}
}
