package cmd

import (
	"database/sql"
	"fmt"
	"strings"

	"cloud.google.com/go/civil"
	_ "github.com/mattn/go-sqlite3"
)

type SummaryRequest struct {
	Start    civil.Date
	End      civil.Date
	Year     int
	Limit    int
	PageSize int
	Page     int
}

type SummaryResponse struct {
	Success bool
	Error   error
	Result  *sql.Rows
}

// Prints the sum of expenses each month
func SummarizeExpenses(sumReq *SummaryRequest) *SummaryResponse {
	err := verifyDatabase()
	if err != nil {
		return &SummaryResponse{
			Success: false,
			Error:   fmt.Errorf("error verifying database: %w", err),
		}
	}

	db, err := connectDB()
	if err != nil {
		return &SummaryResponse{
			Success: false,
			Error:   fmt.Errorf("error connecting to database: %w", err),
		}
	}
	defer db.Close()

	return summaryExec(db, sumReq)
}

func summaryExec(db *sql.DB, req *SummaryRequest) *SummaryResponse {
	connector := "WHERE"
	var sb strings.Builder
	sb.WriteString("SELECT strftime('%Y-%m', date_spent) AS month, sum(amt) AS total_spent FROM expenses")
	if !req.Start.IsZero() {
		sb.WriteString(fmt.Sprintf(" %s date_spent >= '%s'", connector, req.Start.String()))
		connector = "AND"
	}
	if !req.End.IsZero() {
		sb.WriteString(fmt.Sprintf(" %s date_spent <= '%s'", connector, req.End.String()))
		connector = "AND"
	}
	if req.Year != 0 {
		sb.WriteString(fmt.Sprintf(" %s strftime('%%Y', date_spent) = '%d'", connector, req.Year))
		connector = "AND"
	}
	sb.WriteString(" GROUP BY month ORDER BY month")
	if req.Limit != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", req.Limit))
	}
	if req.PageSize != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", req.PageSize))
		if req.Page != 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET %d", req.PageSize*(req.Page-1)))
		}
	}

	sumQuery := sb.String()

	rows, err := db.Query(sumQuery)
	if err != nil {
		return &SummaryResponse{
			Success: false,
			Error:   fmt.Errorf("error calculating summary: %w", err),
		}
	}

	return &SummaryResponse{
		Success: true,
		Result:  rows,
	}
}
