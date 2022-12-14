package repository

import "database/sql"

type Product interface{}

type ProductRepository struct {
	db *sql.DB
}

func newProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}
