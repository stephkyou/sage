package cmd

import (
	"database/sql"
	"fmt"
)

type CountRequest struct {
	Type string
}

type CountResponse struct {
	Success bool
	Error   error
	Result  int
}

// CountExpenses retrieves the total number of expenses
func CountExpenses(db *sql.DB, req *CountRequest) *CountResponse {
	var count int
	if req.Type == "log" {
		rows, err := db.Query("SELECT COUNT(*) FROM expenses")
		if err != nil {
			return &CountResponse{
				Success: false,
				Error:   fmt.Errorf("error querying expenses: %w", err),
			}
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				return &CountResponse{
					Success: false,
					Error:   fmt.Errorf("error scanning count: %w", err),
				}
			}
		}
	} else if req.Type == "summary" {
		rows, err := db.Query("SELECT COUNT(*) FROM expenses GROUP BY strftime('%Y-%m', date_spent)")
		if err != nil {
			return &CountResponse{
				Success: false,
				Error:   fmt.Errorf("error querying expenses: %w", err),
			}
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				return &CountResponse{
					Success: false,
					Error:   fmt.Errorf("error scanning count: %w", err),
				}
			}
		}
	} else {
		return &CountResponse{
			Success: false,
			Error:   fmt.Errorf("invalid request type: %s", req.Type),
		}
	}
	return &CountResponse{Success: true, Result: count}
}
