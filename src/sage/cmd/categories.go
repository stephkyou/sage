package cmd

import (
	"database/sql"
	"fmt"
)

type CategoryRequest struct {
	Subcommand      string
	CategoryName    string
	NewCategoryName string
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
	} else if req.Subcommand == "delete" {
		return deleteCategory(db, req)
	} else if req.Subcommand == "edit" {
		return editCategory(db, req)
	} else {
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

// deleteCategory removes a category from the database
func deleteCategory(db *sql.DB, req *CategoryRequest) *CategoryResponse {
	rows, err := db.Query(fmt.Sprintf("SELECT name FROM categories WHERE name = '%s'", req.CategoryName))
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error querying 'categories' table: %w", err),
		}
	}
	if !rows.Next() {
		rows.Close()
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("category '%s' not found", req.CategoryName),
		}
	}

	rows.Close()

	_, err = db.Exec(fmt.Sprintf("DELETE FROM categories WHERE name = '%s'", req.CategoryName))
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error deleting category from 'categories' table: %w", err),
		}
	}

	return &CategoryResponse{
		Success:    true,
		Subcommand: req.Subcommand,
	}
}

// editCategory changes the name of a category in the database
func editCategory(db *sql.DB, req *CategoryRequest) *CategoryResponse {
	rows, err := db.Query(fmt.Sprintf("SELECT name FROM categories WHERE name = '%s'", req.CategoryName))
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error querying 'categories' table: %w", err),
		}
	}
	if !rows.Next() {
		rows.Close()
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("category '%s' not found", req.CategoryName),
		}
	}

	rows.Close()

	txn, err := db.Begin()
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error starting transaction: %w", err),
		}
	}

	defer func() {
		_ = txn.Rollback()
	}()

	_, err = db.Exec(fmt.Sprintf("INSERT INTO categories (name) VALUES ('%s')", req.NewCategoryName))
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error adding new category to 'categories' table: %w", err),
		}
	}
	_, err = db.Exec(fmt.Sprintf("UPDATE expenses SET category = '%s' WHERE category = '%s'", req.NewCategoryName, req.CategoryName))
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error updating category in 'expenses' table: %w", err),
		}
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM categories WHERE name = '%s'", req.CategoryName))
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error updating category in 'categories' table: %w", err),
		}
	}

	err = txn.Commit()
	if err != nil {
		return &CategoryResponse{
			Success: false,
			Error:   fmt.Errorf("error committing transaction: %w", err),
		}
	}

	return &CategoryResponse{
		Success:    true,
		Subcommand: req.Subcommand,
	}
}
