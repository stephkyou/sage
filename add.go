package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type AddRequest struct {
	Expense Expense
}

// AddExpense adds an expense to the database
func AddExpense(addReq *AddRequest) int {
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

	// insert expense
	_, err = db.Exec(fmt.Sprintf("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('%s', '%s', '%s', %f)",
		addReq.Expense.Date.String(),
		addReq.Expense.Location,
		addReq.Expense.Description,
		addReq.Expense.Amount.AsMajorUnits()))
	if err != nil {
		log.Println("error inserting expense into 'expenses' table: " + err.Error())
		return 1
	}

	return 0
}
