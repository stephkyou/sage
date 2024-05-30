package data

import (
	"cloud.google.com/go/civil"
	"github.com/Rhymond/go-money"
)

type Expense struct {
	Id          int          `json:"-"`
	Date        civil.Date   `json:"date"`
	Location    string       `json:"location,omitempty"`
	Description string       `json:"description,omitempty"`
	Amount      *money.Money `json:"amount"`
}
