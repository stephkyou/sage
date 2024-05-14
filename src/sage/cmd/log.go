package cmd

import (
	"database/sql"
	"fmt"
	"strings"

	"cloud.google.com/go/civil"
	_ "github.com/mattn/go-sqlite3"
)

type LogRequest struct {
	Start    civil.Date
	End      civil.Date
	Year     int
	Month    int
	Limit    int
	PageSize int
	Page     int
	ShowId   bool
}

type LogResponse struct {
	Success bool
	Error   error
	ShowId  bool
	Result  *sql.Rows
}

// LogExpenses retrieves the list of expenses corresponding to the given options and returns the date, location,
// description, and amount (and optionally the expense ID)
func LogExpenses(db *sql.DB, req *LogRequest) *LogResponse {
	connector := "WHERE"
	var sb strings.Builder
	sb.WriteString("SELECT ")
	if req.ShowId {
		sb.WriteString("id, ")
	}
	sb.WriteString("date_spent, location, description, amt FROM expenses")
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
	if req.Month != 0 {
		sb.WriteString(fmt.Sprintf(" %s strftime('%%m', date_spent) = '%d'", connector, req.Month))
		connector = "AND"
	}
	sb.WriteString(" ORDER BY date_spent, id")
	if req.Limit != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", req.Limit))
	}
	if req.PageSize != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", req.PageSize))
		if req.Page != 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET %d", req.PageSize*(req.Page-1)))
		}
	}
	logQuery := sb.String()

	rows, err := db.Query(logQuery)
	if err != nil {
		return &LogResponse{
			Success: false,
			Error:   fmt.Errorf("error retrieving expenses: %w", err),
		}
	}

	return &LogResponse{
		Success: true,
		Result:  rows,
		ShowId:  req.ShowId,
	}
}
