package service

import "shop/internal/repository"

type Product interface{}

type ProductService struct {
	repository repository.Product
}

func newProductService(repository repository.Product) *ProductService {
	return &ProductService{
		repository: repository,
	}
}
