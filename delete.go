package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DeleteRequest struct {
	Id string
}

// DeleteExpense removes an expense from the database
func DeleteExpense(delReq *DeleteRequest) int {
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

	// delete expense
	_, err = db.Exec(fmt.Sprintf("DELETE FROM expenses WHERE id = '%s'", delReq.Id))
	if err != nil {
		log.Println("error deleting expense from 'expenses' table: " + err.Error())
		return 1
	}

	return 0
}
