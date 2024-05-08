package cmd

import (
	"database/sql"
	"errors"
	"io/fs"
	"os"
)

var createTableQuery string = `CREATE TABLE IF NOT EXISTS expenses (
	id INTEGER PRIMARY KEY,
	date_spent DATE NOT NULL,
	location VARCHAR(255),
	description VARCHAR(255),
	amt DECIMAL(19,4) NOT NULL
	)`

// verifyDatabase checks if the sage folder and sage.db database exists. Creates the necessary folder and SQLite file
// if it doesn't.
func verifyDatabase() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return errors.New("error getting user home directory: " + err.Error())
	}

	if _, err := os.Stat(dirname + "/sage"); errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(dirname+"/sage", 0755)
		if err != nil {
			return errors.New("error creating sage directory: " + err.Error())
		}
	}

	if _, err := os.Stat(dirname + "/sage/sage.db"); errors.Is(err, fs.ErrNotExist) {
		file, err := os.Create(dirname + "/sage/sage.db")
		if err != nil {
			return errors.New("error creating database file: " + err.Error())
		}
		file.Close()
	}

	return nil
}

// connectDB connects to the SQLite database and initializes the `expenses` table if it doesn't exist.
func connectDB() (*sql.DB, error) {
	// Connect to database
	db, err := sql.Open("sqlite3", "sage.db")
	if err != nil {
		return nil, errors.New("error connecting to database: " + err.Error())
	}

	// create `expenses` table if it doesn't exist
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, errors.New("error initializing 'expenses' table: " + err.Error())
	}

	return db, nil
}
