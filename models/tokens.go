package models

import "time"

type Token struct {
	ID         int
	SellerID   int
	Signignkey string
	Token      string
	date       time.Time
}
