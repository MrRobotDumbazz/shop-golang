package service 

type Service struct {
	Seller 
	Product 
	Auth
}
func NewServices(repositories *repository.Repository) *Service {
	return &Service {
		Seller: newSellerService(repositories.Seller),
		Product: newProductService(repositories.Product),
		Auth: newAuthService(repositories.Auth)
	}
}