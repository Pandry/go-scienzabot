package database

import "database/sql"

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The types.go file contains only the SQLiteDB embedded type, that is a sql.DB pointer.
// This was made to be sure to take in input only a SQLite database

//SQLiteDB is an instance of the SQL database
type SQLiteDB struct {
	*sql.DB
}
