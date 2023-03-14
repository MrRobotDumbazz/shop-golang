package models

type Product struct {
	ID          int
	SellerID    int
	Name        string
	Company     string
	Description string
	Price       float64
	Category    string
}
