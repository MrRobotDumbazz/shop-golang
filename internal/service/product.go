package service

import "shop/internal/repository"

//go:generate mockgen -source=product.go -destination=mocks/mock.go

type Product interface{}

type ProductService struct {
	repository repository.Product
}

func newProductService(repository repository.Product) *ProductService {
	return &ProductService{
		repository: repository,
	}
}
