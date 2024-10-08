package cmd

import (
	"database/sql"
	"errors"
	"io/fs"
	"os"
)

const (
	CREATE_TABLE_QUERY string = `CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY,
		date_spent DATE NOT NULL,
		location VARCHAR(255),
		description VARCHAR(255),
		category VARCHAR(255),
		amt INTEGER NOT NULL,
		FOREIGN KEY (category) REFERENCES categories(name)
		)`
	CREATE_CATEGORY_TABLE_QUERY string = `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL
		)`
	SAGE_DB_NAME string = "sage.db"
	TEST_DB_NAME string = "test.db"
)

// ConnectDB connects to the given database, or creates it if it doesn't exist. Also initializes the `expenses` table
// if it doesn't exist.
func ConnectDB(db_name string) (*sql.DB, error) {
	// Verify database
	path, err := verifyDatabase(db_name)
	if err != nil {
		return nil, errors.New("error verifying database: " + err.Error())
	}

	// Connect to database
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		return nil, errors.New("error connecting to database: " + err.Error())
	}

	// create `expenses` table if it doesn't exist
	_, err = db.Exec(CREATE_TABLE_QUERY)
	if err != nil {
		return nil, errors.New("error initializing 'expenses' table: " + err.Error())
	}

	// create `categories` table if it doesn't exist
	_, err = db.Exec(CREATE_CATEGORY_TABLE_QUERY)
	if err != nil {
		return nil, errors.New("error initializing 'categories' table: " + err.Error())
	}

	return db, nil
}

// verifyDatabase checks if the sage folder and given database exists. Creates the necessary folder and SQLite file
// if it doesn't. returns the path to the database file.
func verifyDatabase(db_name string) (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("error getting user home directory: " + err.Error())
	}

	if _, err := os.Stat(dirname + "/sage"); errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(dirname+"/sage", 0755)
		if err != nil {
			return "", errors.New("error creating sage directory: " + err.Error())
		}
	}

	if _, err := os.Stat(dirname + "/sage/" + db_name); errors.Is(err, fs.ErrNotExist) {
		file, err := os.Create(dirname + "/sage/" + db_name)
		if err != nil {
			return "", errors.New("error creating database file: " + err.Error())
		}
		file.Close()
	}

	return dirname + "/sage/" + db_name, nil
}
