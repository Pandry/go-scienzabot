package database

import (
	"database/sql"
	"errors"
	"log"
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

func (db *SQLiteDB) selectMultipleQuery() ([]Bookmark, error) {
	rows, err := db.Query("SELECT `ID`, `UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent` FROM Bookmarks")
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	bkms := make([]Bookmark, 0)
	for rows.Next() {
		var (
			id, userID, groupID, messageID, status int64
			messageContent, alias                  string
		)
		if err = rows.Scan(&id, &userID, &groupID, &messageID, &alias, &status, &messageContent); err != nil {
			db.AddLogEvent(Log{Event: "_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			bkms = append(bkms, Bookmark{ID: id, UserID: userID, GroupID: groupID, MessageID: messageID, Alias: alias, Status: status, MessageContent: messageContent})
		}
	}
	if !rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "_RowsNotFetched", Message: "Some rows in the query were not fetched", Error: err.Error()})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return bkms, err
}

/*























 */

func (db *SQLiteDB) selectMultipleQueryOlddd() error {
	rows, err := db.Query("SELECT * FROM Table WHERE `a`=?", "a")
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	strings := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			db.AddLogEvent(Log{Event: "_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		}
		strings = append(strings, name)
	}
	if !rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "_RowNotFetched", Message: "Some rows in the query were not fetched", Error: err.Error()})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return err
}
