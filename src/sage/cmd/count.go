package cmd

import (
	"database/sql"
	"fmt"
)

type CountRequest struct {
	Query string
}

type CountResponse struct {
	Success bool
	Error   error
	Result  int
}

// CountExpenses retrieves the total number of expenses
func CountExpenses(db *sql.DB, req *CountRequest) *CountResponse {
	rows, err := db.Query("SELECT COUNT(*) FROM expenses")
	if err != nil {
		return &CountResponse{
			Success: false,
			Error:   fmt.Errorf("error querying expenses: %w", err),
		}
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return &CountResponse{
				Success: false,
				Error:   fmt.Errorf("error scanning count: %w", err),
			}
		}
	}
	return &CountResponse{Success: true, Result: count}
}
