package cmd

import (
	"database/sql"
	"fmt"
)

type CategoryRequest struct {
	Subcommand   string
	CategoryName string
}

type CategoryResponse struct {
	Success    bool
	Error      error
	Subcommand string
	Result     *sql.Rows
}

// ExpenseCategory retrieves the list of categories and returns their names
func ExpenseCategory(db *sql.DB, req *CategoryRequest) *CategoryResponse {
	if req.Subcommand == "add" {
		return addCategory(db, req)
	}

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

// addCategory adds a new category to the database
func addCategory(db *sql.DB, req *CategoryRequest) *CategoryResponse {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO categories (name) VALUES ('%s')", req.CategoryName))
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error adding category to 'categories' table: %w", err),
		}
	}

	return &CategoryResponse{
		Success:    true,
		Subcommand: req.Subcommand,
	}
}
