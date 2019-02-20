package database

import (
	"database/sql"
	"log"

	"scienzabot/consts"

	_ "github.com/mattn/go-sqlite3" //apk add --update gcc musl-dev
)

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
