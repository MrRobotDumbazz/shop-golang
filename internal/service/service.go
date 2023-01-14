package service

import "shop/internal/repository"

type Service struct {
	Seller
	Product
	Auth
}

func NewServices(repositories *repository.Repository, redis repository.RedisRepository) *Service {
	return &Service{
		Seller:  newSellerService(repositories.Seller),
		Product: newProductService(repositories.Product),
		Auth:    newAuthService(repositories.Auth, redis),
	}
}
