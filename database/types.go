package database

import "database/sql"

//SQLiteDB is an instance of the SQL database
type SQLiteDB struct {
	*sql.DB
}
