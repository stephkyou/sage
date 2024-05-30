package cmd

import (
	"database/sql"
	"fmt"
	"sage/src/sage/data"

	_ "github.com/mattn/go-sqlite3"
)

type AddRequest struct {
	Expense data.Expense
}

type AddResponse struct {
	Success bool
	Error   error
}

// addExpense adds an expense to the database
func AddExpense(db *sql.DB, req *AddRequest) *AddResponse {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('%s', '%s', '%s', %d)",
		req.Expense.Date.String(),
		req.Expense.Location,
		req.Expense.Description,
		req.Expense.Amount.Amount()))
	if err != nil {
		return &AddResponse{
			Success: false,
			Error:   fmt.Errorf("error adding expense to 'expenses' table: %w", err),
		}
	}

	return &AddResponse{Success: true}
}
