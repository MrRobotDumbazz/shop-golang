package repository

import (
	"database/sql"
	"shop/models"
)

type Product interface {
	GetNewAllProducts() ([]models.Product, error)
	GetProductByProductID(id int) (*models.Product, error)
	CreateProduct(p *models.Product) error
	GetProductByCategory(category string) ([]models.Product, error)
}

type ProductRepository struct {
	db *sql.DB
}

func newProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetNewAllProducts() ([]models.Product, error) {
	var products []models.Product
	query := "SELECT *  FROM shodb.products ORDER BY DESC;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		product := models.Product{}
		if err := rows.Scan(&product.ID, &product.SellerID, &product.Name, &product.Company, &product.Description, &product.Price, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *ProductRepository) GetProductByProductID(id int) (*models.Product, error) {
	p := &models.Product{}
	err := r.db.QueryRow("SELECT * FROM shopdb.products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Category)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) GetProductByCategory(category string) ([]models.Product, error) {
	var products []models.Product
	query := "SELECT *  FROM shopdb.products ORDER BY DESC WHERE category = ?;"
	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		product := models.Product{}
		if err := rows.Scan(&product.ID, &product.SellerID, &product.Name, &product.Company, &product.Description, &product.Price, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *ProductRepository) CreateProduct(p *models.Product) error {
	if _, err := r.db.Exec("INSERT INTO shopdb.products (seller_id, name_product, company, description, category, price) VALUES (?, ?, ?, ?, ?, ?)", p.SellerID, p.Name, p.Company, p.Description, p.Category, p.Price); err != nil {
		return err
	}
	return nil
}
