package database

import (
	"database/sql"
	"errors"
)

func (db *SQLiteDB) insertQuery() error {
	query, err := db.Exec("INSERT INTO Table (`a`, `b`, `c`) VALUES (?,?,?)", "")
	if err != nil {
		db.AddLogEvent(Log{Event: "_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

func (db *SQLiteDB) deleteQuery() error {
	query, err := db.Exec("DELETE FROM Table WHERE a=?", "a")
	if err != nil {
		db.AddLogEvent(Log{Event: "_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

func (db *SQLiteDB) updateQuery() error {
	query, err := db.Exec("UPDATE Table SET `Field`=?, `Field2`=? WHERE `ID`=?", "a", "", "")
	if err != nil {
		db.AddLogEvent(Log{Event: "_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

func (db *SQLiteDB) selectSingleQuery() error {
	var (
		name string
		Age  sql.NullInt64
	)
	err := db.QueryRow("SELECT * FROM Table WHERE `ID`=?", "a").Scan(&name, &Age)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return err
	case err != nil:
		db.AddLogEvent(Log{Event: "_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return err
	default:
		return nil
	}
}

func (db *SQLiteDB) selectMultipleQuery() error {
	query, err := db.Query("SELECT * FROM Table WHERE `a`=?", "a")
	defer query.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	strings := make([]string, 0)
	for query.Next() {
		var name string
		if err := query.Scan(&name); err != nil {
			db.AddLogEvent(Log{Event: "_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		}
		strings = append(strings, name)
	}
	if !query.NextResultSet() {
		db.AddLogEvent(Log{Event: "_RowNotFetched", Message: "Some rows in the query were not fetched", Error: err.Error()})

	}

	return err
}
