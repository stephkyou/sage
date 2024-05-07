package main

import (
	"database/sql"
	"fmt"
	"log"
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

// Prints the sum of expenses each month
func SummarizeExpenses(sumReq *SummaryRequest) int {
	// Connect to database
	db, err := sql.Open("sqlite3", "sage.db")
	if err != nil {
		fmt.Println("error connecting to database:", err)
		return 1
	}
	defer db.Close()

	// create `expenses` table if it doesn't exist
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Println("error initializing 'expenses' table: " + err.Error())
		return 1
	}

	// construct query
	connector := "WHERE"
	var sb strings.Builder
	sb.WriteString("SELECT strftime('%Y-%m', date_spent) AS month, sum(amt) AS total_spent FROM expenses")
	if !sumReq.Start.IsZero() {
		sb.WriteString(fmt.Sprintf(" %s date_spent >= '%s'", connector, sumReq.Start.String()))
		connector = "AND"
	}
	if !sumReq.End.IsZero() {
		sb.WriteString(fmt.Sprintf(" %s date_spent <= '%s'", connector, sumReq.End.String()))
		connector = "AND"
	}
	if sumReq.Year != 0 {
		sb.WriteString(fmt.Sprintf(" %s strftime('%%Y', date_spent) = '%d'", connector, sumReq.Year))
		connector = "AND"
	}
	if sumReq.Limit != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", sumReq.Limit))
	}
	if sumReq.PageSize != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", sumReq.PageSize))
		if sumReq.Page != 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET %d", sumReq.PageSize*(sumReq.Page-1)))
		}
	}
	sb.WriteString(" GROUP BY month ORDER BY month")

	sumQuery := sb.String()

	rows, err := db.Query(sumQuery)
	if err != nil {
		log.Println("error calculating summary: " + err.Error())
		return 1
	}

	for rows.Next() {
		var month string
		var totalSpent float64
		err = rows.Scan(&month, &totalSpent)
		if err != nil {
			log.Println("error reading calculated summary: " + err.Error())
		}
		fmt.Printf("%s: $%.2f\n", month, totalSpent)
	}

	return 0
}
