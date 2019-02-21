package database

import (
	"database/sql"
	"errors"
	"scienzabot/consts"
)

/*
CREATE TABLE IF NOT EXISTS 'BotStrings' (
	'ID'		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'		TEXT NOT NULL UNIQUE,
	'Value'		TEXT DEFAULT 'Not implemented',
	'Locale'	TEXT DEFAULT '` + DefaultLocale + `',
	CONSTRAINT con_botstrings_key_locale_unique UNIQUE ('Key','Locale')
);
*/

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

//GetBotStringValueOrDefault returns the value in the user's locale or in the default one
func (db *SQLiteDB) GetBotStringValueOrDefault(key string, locale string) (string, error) {
	if db.BotStringExists(key, locale) {
		return db.GetBotStringValue(key, locale)
	}
	return db.GetBotStringValue(key, consts.DefaultLocale)
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
		db.AddLogEvent(Log{Event: "SetBotStringValue_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}
