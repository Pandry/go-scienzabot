package database

import "errors"

func (db *SQLiteDB) insertQuery() error {
	stmt, err := db.Prepare("INSERT INTO VALUES (?,?,?)")
	if err != nil {
		db.AddLogEvent(Log{Event: "_QueryFailed", Message: "Impossible to create the  preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		db.AddLogEvent(Log{Event: "_ExecutionQueryFailed", Message: "Impossible to execute the  query", Error: err.Error()})
		return err
	}

	rows, err := res.RowsAffected()

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
	stmt, err := db.Prepare("DELETE FROM Table WHERE `ID=?`")
	if err != nil {
		db.AddLogEvent(Log{Event: "_QueryFailed", Message: "Impossible to create the  preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		db.AddLogEvent(Log{Event: "_ExecutionQueryFailed", Message: "Impossible to execute the  query", Error: err.Error()})
		return err
	}

	rows, err := res.RowsAffected()

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
	stmt, err := db.Prepare("UPDATE Table SET `Field`=?, `Field2`=? WHERE `ID`=?")
	if err != nil {
		db.AddLogEvent(Log{Event: "_QueryFailed", Message: "Impossible to create the  preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		db.AddLogEvent(Log{Event: "_ExecutionQueryFailed", Message: "Impossible to execute the  query", Error: err.Error()})
		return err
	}

	rows, err := res.RowsAffected()

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

func (db *SQLiteDB) selectMultipleQuery() error {
	stmt, err := db.Prepare("SELECT * FROM Table WHERE `ID`=?")
	if err != nil {
		db.AddLogEvent(Log{Event: "_QueryFailed", Message: "Impossible to create the  preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	if err != nil {
		db.AddLogEvent(Log{Event: "_ExecutionQueryFailed", Message: "Impossible to execute the  query", Error: err.Error()})
		return err
	}

	err = res.Err()

	if err != nil {
		db.AddLogEvent(Log{Event: "_SelectQueryFailed", Message: "Impossible to execute the query: an error verified", Error: err.Error()})
		return err
	}

	strings := make([]string, 0)
	for res.Next() {
		var name string

		if err := res.Scan(&name); err != nil {
			db.AddLogEvent(Log{Event: "_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		}

		strings = append(strings, name)

	}

	return err
}
