package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Connection struct {
	db *sql.DB
}

var dbCon Connection

// Create a new connection to the database
// Creates the database file if it doesn't exist and creates the initial table if it doesn't exist
func CreateConnection() (Connection, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return dbCon, fmt.Errorf("Failed to open or create the database file: %s", err)
	}
	dbCon = Connection{db}

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS task_item (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT,
			Description TEXT,
			DueDate DATETIME,
			PriorityLevel INT,
			Completed TINYINT
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return dbCon, fmt.Errorf("Failed to create the initial table: %s", err)
	}

	return dbCon, nil
}

// Execute a query and return the result as a sql.Rows object
func (dbCon Connection) query(query string, bindings ...any) (*sql.Rows, error) {
	return dbCon.db.Query(query, bindings...)
}

// Execute a statement and return the result as a sql.Result object
func (dbCon Connection) executeStatement(query string, bindings ...any) (sql.Result, error) {
	tx, err := dbCon.db.Begin()
	if err != nil {
		return nil, err

	}

	defer tx.Commit()
	result, err := tx.Exec(query, bindings...)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Close the connection to the database
func (dbCon Connection) CloseClient() {
	dbCon.db.Close()
}
