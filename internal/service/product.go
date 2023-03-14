package service

import (
	"errors"
	"log"
	"shop/internal/repository"
	"shop/models"
	"strings"
)

//go:generate mockgen -source=product.go -destination=mocks/mock.go
var (
	ErrInvalidQueryRequest = errors.New("invalid query request")
)

type Product interface {
	GetAllProductsBy(query map[string][]string) ([]models.Product, error)
	GetProductByProductID(id int) (*models.Product, error)
	CreateProduct(p *models.Product) error
	GetNewAllProducts() ([]models.Product, error)
}

type ProductService struct {
	repository repository.Product
}

func newProductService(repository repository.Product) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetAllProductsBy(query map[string][]string) ([]models.Product, error) {
	var (
		products []models.Product
		err      error
	)
	for key, val := range query {
		switch key {
		case "category":
			products, err = s.repository.GetProductByCategory(strings.Join(val, ""))
			if err != nil {
				log.Println("error:products:GetProductByCategory ", err)
				return nil, err
			}
		default:
			log.Println("error:post:GetAllProductBy:default ", err)
			return nil, ErrInvalidQueryRequest
		}
	}
	return products, nil
}

func (s *ProductService) GetProductByProductID(id int) (*models.Product, error) {
	product, err := s.repository.GetProductByProductID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) CreateProduct(p *models.Product) error {
	if err := s.repository.CreateProduct(p); err != nil {
		return err
	}
	return nil
}

func (s *ProductService) GetNewAllProducts() ([]models.Product, error) {
	products, err := s.repository.GetNewAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}
