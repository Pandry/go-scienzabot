package database

import (
	"database/sql"
	"errors"
)

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
		db.AddLogEvent(Log{Event: "SetBotSettingValue_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
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
