package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

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

// LogExpenses retrieves the list of expenses corresponding to the given options and displays the date, location,
// description, and amount (and optionally the expense ID)
func LogExpenses(logReq *LogRequest) int {
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
	sb.WriteString("SELECT ")
	if logReq.ShowId {
		sb.WriteString("id, ")
	}
	sb.WriteString("date_spent, location, description, amt FROM expenses")
	if !logReq.Start.IsZero() {
		sb.WriteString(fmt.Sprintf(" %s date_spent >= '%s'", connector, logReq.Start.String()))
		connector = "AND"
	}
	if !logReq.End.IsZero() {
		sb.WriteString(fmt.Sprintf(" %s date_spent <= '%s'", connector, logReq.End.String()))
		connector = "AND"
	}
	if logReq.Year != 0 {
		sb.WriteString(fmt.Sprintf(" %s strftime('%%Y', date_spent) = '%d'", connector, logReq.Year))
		connector = "AND"
	}
	if logReq.Month != 0 {
		sb.WriteString(fmt.Sprintf(" %s strftime('%%m', date_spent) = '%d'", connector, logReq.Month))
		connector = "AND"
	}
	sb.WriteString(" ORDER BY date_spent, id")
	if logReq.Limit != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", logReq.Limit))
	}
	if logReq.PageSize != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", logReq.PageSize))
		if logReq.Page != 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET %d", logReq.PageSize*(logReq.Page-1)))
		}
	}
	logQuery := sb.String()

	rows, err := db.Query(logQuery)
	if err != nil {
		log.Println("error retrieving expenses: " + err.Error())
		return 1
	}
	defer rows.Close()

	for rows.Next() {
		var date time.Time
		var location string
		var description string
		var amt float64
		if logReq.ShowId {
			var id int
			err := rows.Scan(&id, &date, &location, &description, &amt)
			if err != nil {
				log.Println("error reading retrieved expenses: " + err.Error())
				return 1
			}
			fmt.Printf("%d | %s | %s | %s | $%.2f\n", id, date.Format("2006-01-02"), location, description, amt)
		} else {
			err := rows.Scan(&date, &location, &description, &amt)
			if err != nil {
				log.Println("error reading retrieved expenses: " + err.Error())
				return 1
			}
			fmt.Printf("%s | %s | %s | $%.2f\n", date.Format("2006-01-02"), location, description, amt)
		}
	}

	return 0
}
