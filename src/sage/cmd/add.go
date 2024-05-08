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

// AddExpense adds an expense to the database
func AddExpense(addReq *AddRequest) *AddResponse {
	err := verifyDatabase()
	if err != nil {
		return &AddResponse{
			Success: false,
			Error:   fmt.Errorf("error verifying database: %w", err),
		}
	}

	db, err := connectDB()
	if err != nil {
		return &AddResponse{
			Success: false,
			Error:   fmt.Errorf("error connecting to database: %w", err),
		}
	}
	defer db.Close()

	return addExec(db, addReq)
}

// addExec adds an expense to the database
func addExec(db *sql.DB, req *AddRequest) *AddResponse {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('%s', '%s', '%s', %f)",
		req.Expense.Date.String(),
		req.Expense.Location,
		req.Expense.Description,
		req.Expense.Amount.AsMajorUnits()))
	if err != nil {
		return &AddResponse{
			Success: false,
			Error:   fmt.Errorf("error adding expense to 'expenses' table: %w", err),
		}
	}

	return &AddResponse{Success: true}
}
