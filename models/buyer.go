package models

type Client struct {
	ID       int    `json:"id"`
	Email    string `json:"username"`
	Password string `json:"password"`
}
