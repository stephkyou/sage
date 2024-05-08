package data

import (
	"cloud.google.com/go/civil"
	"github.com/Rhymond/go-money"
)

type Expense struct {
	Id          int
	Date        civil.Date
	Location    string
	Description string
	Amount      *money.Money
}
