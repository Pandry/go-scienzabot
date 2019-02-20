package database

import (
	"database/sql"
	"errors"
)

//GetStringValue searches for the string value in the database
func (db *SQLiteDB) GetStringValue(key string, group int, locale string) (string, error) {
	res := ""
	err := db.QueryRow("SELECT Value FROM `Strings` WHERE `Key` = ? AND `Group` = ? AND `Locale` = ?",
		key, group, locale).Scan(&res)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetStringValue_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return res, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetStringValue_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return res, err
	default:
		return res, err
	}
}

//SetStringValue sets a value in the bot settings table
func (db *SQLiteDB) SetStringValue(key string, value string, group int, locale string) error {
	query, err := db.Exec(
		"INSERT INTO Settings (`Key`, `Value` , `Group`) VALUES (?,?,?) "+
			"ON CONFLICT(`Key`, `Group`) DO UPDATE SET `Value` = Excluded.Value",
		key, value, group)
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
		db.AddLogEvent(Log{Event: "SetStringValue_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
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
	stmt, err := db.Prepare("SELECT Value FROM `Strings` WHERE `Key` = ? AND `Group` = ? AND `Locale` = ? ")
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
	stmt, err := db.Prepare("SELECT Value FROM `Strings` WHERE `Key` = ? AND `Group` = ? AND `Locale` = ? ")
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
	stmt, err := db.Prepare("INSERT INTO Settings (`Key`, `Value` , `Group`) VALUES (?,?,?) ON CONFLICT(`Key`, `Group`) DO UPDATE SET `Value` = Excluded.Value")
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
