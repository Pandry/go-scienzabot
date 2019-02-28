package database

import (
	"database/sql"
	"log"

	"scienzabot/consts"

	_ "github.com/mattn/go-sqlite3" //apk add --update gcc musl-dev
)

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The database.go file contains the initialization function and return an instance of the
//	database

//InitDatabaseConnection occupied of establishing the first connection to the database
func InitDatabaseConnection(dbPath string) (sqlitedbStruct *SQLiteDB, err error) {
	//Initializes the connection to the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	//Check for errors
	if err != nil {
		log.Fatal(err)
	}

	//Execute the initialization query
	_, err = db.Exec(consts.InitSQLString)
	if err != nil {
		return nil, err
	}

	dbs := &SQLiteDB{db}

	return dbs, nil
}
