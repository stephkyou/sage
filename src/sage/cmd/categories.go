package cmd

import (
	"database/sql"
	"fmt"
)

type CategoryRequest struct {
	Subcommand string
}

type CategoryResponse struct {
	Success bool
	Error   error
	Result  *sql.Rows
}

// ExpenseCategory retrieves the list of categories and returns their names
func ExpenseCategory(db *sql.DB, req *CategoryRequest) *CategoryResponse {
	result, err := db.Query("SELECT name FROM categories")
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error querying 'categories' table: %w", err),
		}
	}

	return &CategoryResponse{
		Success: true,
		Result:  result,
	}
}
