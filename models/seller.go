package models

type Seller struct {
	ID       int    `json:"id"`
	Email    string `json:"username"`
	Password string `json:"password,omitempty"`
	HasToken bool
	Products Product
}
