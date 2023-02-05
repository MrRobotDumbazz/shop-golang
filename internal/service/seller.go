package service

import "shop/internal/repository"

//go:generate mockgen -source=seller.go -destination=mocks/mockseller.go

type Seller interface{}

type SellerService struct {
	repository repository.Seller
}

func newSellerService(repository repository.Seller) *SellerService {
	return &SellerService{
		repository: repository,
	}
}
