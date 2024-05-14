package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DeleteRequest struct {
	Id int
}

type DeleteResponse struct {
	Success bool
	Error   error
}

// DeleteExpense removes an expense from the database
func DeleteExpense(db *sql.DB, req *DeleteRequest) *DeleteResponse {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM expenses WHERE id = '%d'", req.Id))
	if err != nil {
		return &DeleteResponse{
			Success: false,
			Error:   fmt.Errorf("error deleting expense from 'expenses' table: %w", err),
		}
	}

	return &DeleteResponse{Success: true}
}
