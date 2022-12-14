package repository

import "database/sql"

type Seller interface{}

type SellerRepository struct {
	db *sql.DB
}

func newSellerRepository(db *sql.DB) *SellerRepository {
	return &SellerRepository{
		db: db,
	}
}
