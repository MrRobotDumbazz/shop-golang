package service

import "shop/internal/repository"

type Seller interface{}

type SellerService struct {
	repository repository.Seller
}

func newSellerService(repository repository.Seller) *SellerService {
	return &SellerService{
		repository: repository,
	}
}
