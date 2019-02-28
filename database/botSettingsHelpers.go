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

// The botSettingsHelpers.go file focuses on the BotSettings table in the database, that is
//	supposed to contain strings that refer to the bot configuration, like the default locale.

//BotSettingExists returns a values that indicates if the key exists in database
func (db *SQLiteDB) BotSettingExists(key string) bool {
	var dummyval int64
	err := db.QueryRow("SELECT 1 FROM `BotSettings` WHERE `Key` = ?",
		key).Scan(&dummyval)
	switch {
	case err == sql.ErrNoRows:
		//db.AddLogEvent(Log{Event: "_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return false
	case err != nil:
		db.AddLogEvent(Log{Event: "BotSettingsExists_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return false
	default:
		return true
	}
}

//GetBotSettingValue searches for the value in the database
func (db *SQLiteDB) GetBotSettingValue(key string) (string, error) {
	var val string
	err := db.QueryRow("SELECT Value FROM `BotSettings` WHERE Key = ?", key).Scan(&val)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetBotSettingValue_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return "", err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetBotSettingValue_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return "", err
	default:
		return val, nil
	}
	/*
	   THIS COMMENT IS BEING KEPT SINCE IT EXPLAINED HOW THE QUERIES WERE DONE BEFORE A REFACTOR,
	     EXPLAINING THINGS
	   DO NOT REMOVE, PLEASE (At least not before more comments are added to the code above)
	   	//We're prepaing a query to execute later
	   	stmt, err := db.Prepare("SELECT Value FROM `BotSettings` WHERE Key = ?")
	   	if err != nil {
	   		return "", err
	   	}
	   	//We want to close the connection to the database once we stop using it
	   	defer stmt.Close()
	   	//The setting value will go on this string
	   	var val string
	   	//Then we execute the query passing the key to the scan function
	   	err = stmt.QueryRow(key).Scan(&val)
	   	if err != nil {
	   		return "", err
	   	}
	   	//Finally, we return the result
	   	return val, nil
	*/
}

//SetBotSettingValue sets a value in the bot settings table
func (db *SQLiteDB) SetBotSettingValue(key string, value string) error {

	query, err := db.Exec("INSERT INTO BotSettings (Key,Value) VALUES (?, ?) ON CONFLICT(Key) DO UPDATE SET Key = Key;", key, value)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetBotSettingValue_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetBotSettingValue_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetBotSettingValue_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err

	// == OLD CODE ==
	// KEPT ONLY TO SHOW HOW MUCH THE OLD CODE USED TO SUCK
	//
	//
	//here we're starting a transaction with the database.
	//This is becase we are changing the database and we want to write the changes to
	//the file in a permanent way after we're done
	//
	// ACTUALLY the lib automacically starts a transaction, so it is not needed,
	// but I'm keeping the comments anyway
	//
	/*tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	//stmt, err := tx.Prepare("INSERT INTO BotSettings (Key,Value) VALUES (?, ?) ON DUPLICATE KEY UPDATE Value=VALUES(Value)")
	//Then we prepare the query
	//The above statement is how it should be done in a MySQL-like DB using a transaction
	stmt, err := db.Prepare("INSERT INTO BotSettings (Key,Value) VALUES (?, ?) ON CONFLICT(Key) DO UPDATE SET Key = Key;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(key, value)

	if err != nil {
		return err
	}
	//Finally, if everything went good, we commit the database
	//Still, this is old code, but it's the siggested practice
	//err = tx.Commit()
	return err

	*/
}
