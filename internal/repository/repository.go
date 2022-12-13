package repository

import "database/sql"

type Repository struct {
	Client
	Seller
	Auth
	Product
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Client:  newClientRepository(db),
		Seller:  newSellerRepostiroy(db),
		Auth:    newAuthRepository(db),
		Product: newProductRepository(db),
	}
}
