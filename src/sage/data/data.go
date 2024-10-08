package data

import (
	"cloud.google.com/go/civil"
	"github.com/Rhymond/go-money"
)

type Expense struct {
	Id          int          `json:"id"`
	Date        civil.Date   `json:"date"`
	Location    string       `json:"location,omitempty"`
	Description string       `json:"description,omitempty"`
	Category    string       `json:"category,omitempty"`
	Amount      *money.Money `json:"amount"`
}

type Summary struct {
	Month string       `json:"month"`
	Total *money.Money `json:"total"`
}
