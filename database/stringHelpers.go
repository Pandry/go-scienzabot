package database

import (
	"database/sql"
	"errors"
	"scienzabot/consts"
)

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The stringHelpers.go file focuses on the Strings table in the database, that is supposed to
//  have a key-value structure and keep strings from group, like a welcome message for example

//StringExists returns a values that indicates if the key exists in database
func (db *SQLiteDB) StringExists(key string, locale string, group int64) bool {
	if locale == "" {
		locale = consts.DefaultLocale
	}
	var dummyval int64
	err := db.QueryRow("SELECT 1 FROM `Strings` WHERE `Key` = ? AND `Locale` = ? AND `GroupID`=?",
		key, locale, group).Scan(&dummyval)
	switch {
	case err == sql.ErrNoRows:
		//db.AddLogEvent(Log{Event: "_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return false
	case err != nil:
		db.AddLogEvent(Log{Event: "BotStringExists_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return false
	default:
		return true
	}
}

//GetStringValue searches for the string value in the database
func (db *SQLiteDB) GetStringValue(key string, group int64, locale string) (string, error) {
	var res sql.NullString
	if locale == "" {
		locale = consts.DefaultLocale
	}
	err := db.QueryRow("SELECT Value FROM `Strings` WHERE `Key` = ? AND `GroupID` = ? AND `Locale` = ?",
		key, group, locale).Scan(&res)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetStringValue_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return res.String, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetStringValue_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return res.String, err
	default:
		return res.String, err
	}
}

//SetStringValue sets a value in the bot settings table
func (db *SQLiteDB) SetStringValue(key string, value string, group int64, locale string) error {
	if locale == "" {
		locale = consts.DefaultLocale
	}
	query, err := db.Exec(
		"INSERT INTO Strings (`Key`, `Value`, `Locale`, `GroupID`) VALUES (?,?,?,?) "+
			"ON CONFLICT(`Key`, `GroupID`, `Locale`) DO UPDATE SET `Value` = Excluded.Value",
		key, value, locale, group)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetStringValue_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetStringValue_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetStringValue_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

/*

import (
	"database/sql"
	"log"
)

//TODO: return here when basic queries will be done

//For a more detailed explaination about this code, see botSettingsHelpers.go file, in this directory

//String is referred as a string like the welcomemessage, so custom from a group to another

//GetStringValue searches for the string value in the database
func (db *SQLiteDB) GetStringValue(key string, group int, locale string) (string, error) {
	//We're prepaing a query to execute later
	stmt, err := db.Prepare("SELECT Value FROM `Strings` WHERE `Key` = ? AND `GroupID` = ? AND `Locale` = ? ")
	if err != nil {
		return "", err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()
	//The setting value will go on this string
	var val string
	//Then we execute the query passing the key to the scan function
	err = stmt.QueryRow(key, group, locale).Scan(&val)
	if err != nil {
		return "", err
	}
	//Finally, we return the result
	return val, nil
}

//GetDefaultStringValue searches for the string value in the database using first the default
//locale of the group, then the  bot default locale, if nothing found tries without specifying the locale
//and if there's still no result, error is returned
func (db *SQLiteDB) GetDefaultStringValue(key string, group int) (string, error) {

	//We're prepaing a query to execute later
	stmt, err := db.Prepare("SELECT Value FROM `Strings` WHERE `Key` = ? AND `GroupID` = ? AND `Locale` = ? ")
	if err != nil {
		return "", err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()
	//The setting value will go on this string
	var val string
	//Then we execute the query passing the key to the scan function
	err = stmt.QueryRow(key, group, locale).Scan(&val)
	switch {

	case err == sql.ErrNoRows:
		log.Printf("No user with id %d", id)

	case err != nil:
		log.Fatal(err)

	default:
		//Success

	}
	if err != nil {
		return "", err
	}
	//Finally, we return the result
	return val, nil
}

//SetStringValue sets a value in the bot settings table
func (db *SQLiteDB) SetStringValue(key string, value string, group int, locale string) error {
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
}
*/
