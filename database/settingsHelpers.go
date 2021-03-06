package database

import (
	"database/sql"
	"errors"
)

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The settingsHelpers.go file focuses on the Settings table in the database, that is supposed to
//  have a key-value structure and is supposed to keep eventual settings of a group, like an eventual
//  antiflood setting, for example

//SettingExists returns a values that indicates if the key exists in database
func (db *SQLiteDB) SettingExists(key string) bool {
	var dummyval int64
	err := db.QueryRow("SELECT 1 FROM `Settings` WHERE `Key` = ?",
		key).Scan(&dummyval)
	switch {
	case err == sql.ErrNoRows:
		//db.AddLogEvent(Log{Event: "_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return false
	case err != nil:
		db.AddLogEvent(Log{Event: "SettingsExists_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return false
	default:
		return true
	}
}

//GetSettingValue searches for the value in the database
func (db *SQLiteDB) GetSettingValue(key string, group int64) (string, error) {
	val := ""
	err := db.QueryRow("SELECT Value FROM `Settings` WHERE `Key` = ? AND `GroupID` = ?",
		key, group).Scan(&val)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetSettingValue_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return val, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetSettingValue_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return val, err
	default:
		return val, nil
	}
	/*
		//We're prepaing a query to execute later
		stmt, err := db.Prepare("SELECT Value FROM `Settings` WHERE `Key` = ? AND `GroupID` = ? ")
		if err != nil {
			return "", err
		}
		//We want to close the connection to the database once we stop using it
		defer stmt.Close()
		//The setting value will go on this string
		var val string
		//Then we execute the query passing the key to the scan function
		err = stmt.QueryRow(key, group).Scan(&val)
		if err != nil {
			return "", err
		}
		//Finally, we return the result
		return val, nil
	*/
}

//SetSettingValue sets a value in the bot settings table
func (db *SQLiteDB) SetSettingValue(key string, value string, group int64) error {
	query, err := db.Exec("INSERT INTO Settings (`Key`, `Value` , `GroupID`) VALUES (?,?,?) "+
		"ON CONFLICT(`Key`, `GroupID`) DO UPDATE SET `Value` = Excluded.Value", key, value, group)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetSettingValue_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetSettingValue_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetSettingValue_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
	/*
		stmt, err := db.Prepare("INSERT INTO Settings (`Key`, `Value` , `GroupID`) VALUES (?,?,?) ON CONFLICT(`Key`, `GroupID`) DO UPDATE SET `Value` = Excluded.Value")
		if err != nil {
			return err
		}
		defer stmt.Close()

		//And we execute it passing the parameters
		stmt.Exec(key, value, group)

		if err != nil {
			return err
		}

		return err
	*/
}
