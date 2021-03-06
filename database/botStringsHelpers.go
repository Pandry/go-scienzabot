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

// The botStringsHelpers.go file focuses on the BotStrings table in the database, that is
//	supposed to contain immutable strings, valid for each group.
// An example could be the close button text value, or the /info command

//BotStringExists returns a values that indicates if the key exists in database
func (db *SQLiteDB) BotStringExists(key string, locale string) bool {
	if locale == "" {
		locale = consts.DefaultLocale
	}
	var a int64
	err := db.QueryRow("SELECT 1 FROM `BotStrings` WHERE `Key` = ? AND `Locale` = ?",
		key, locale).Scan(&a)
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

//GetBotStringValue searches for the BOT string value in the database
func (db *SQLiteDB) GetBotStringValue(key string, locale string) (string, error) {
	if locale == "" {
		locale = consts.DefaultLocale
	}
	var res sql.NullString
	err := db.QueryRow("SELECT Value FROM `BotStrings` WHERE `Key` = ? AND `Locale` = ?",
		key, locale).Scan(&res)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetBotStringValue_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return res.String, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetBotStringValue_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return res.String, err
	default:
		return res.String, err
	}
}

//GetBotStringValue searches for the BOT string value in the database
func (db *SQLiteDB) getFirstBotStringValue(key string) (string, error) {
	var res sql.NullString
	err := db.QueryRow("SELECT Value FROM `BotStrings` WHERE `Key` = ? LIMIT 1",
		key).Scan(&res)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "getFirstBotStringValue_ErrorNoRows", Message: "Impossible to get rows with key: \"" + key + "\"", Error: err.Error()})
		return res.String, err
	case err != nil:
		db.AddLogEvent(Log{Event: "getFirstBotStringValue_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return res.String, err
	default:
		return res.String, err
	}
}

//GetBotStringValueOrDefaultNoError does not return an error
func (db *SQLiteDB) GetBotStringValueOrDefaultNoError(key string, locale string) string {
	if db.BotStringExists(key, locale) {
		s, _ := db.GetBotStringValue(key, locale)
		return s
	}
	if db.BotStringExists(key, consts.DefaultLocale) {
		s, _ := db.GetBotStringValue(key, consts.DefaultLocale)
		return s
	}
	s, _ := db.getFirstBotStringValue(key)
	return s
}

//GetBotStringValueOrDefault returns the value in the user's locale or in the default one
func (db *SQLiteDB) GetBotStringValueOrDefault(key string, locale string) (string, error) {
	if db.BotStringExists(key, locale) {
		return db.GetBotStringValue(key, locale)
	}
	if db.BotStringExists(key, consts.DefaultLocale) {
		return db.GetBotStringValue(key, consts.DefaultLocale)
	}
	return db.getFirstBotStringValue(key)
}

//SetBotStringValue sets a value in the bot settings table
func (db *SQLiteDB) SetBotStringValue(key string, value string, locale string) error {
	query, err := db.Exec(
		"INSERT INTO BotStrings (`Key`, `Value`, `Locale`) VALUES (?,?,?) "+
			"ON CONFLICT(`Key`, `Locale`) DO UPDATE SET `Value` = Excluded.Value",
		key, value, locale)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetBotStringValue_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetBotStringValue_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetBotStringValue_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}
