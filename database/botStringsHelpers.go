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

//GetBotStringValue searches for the BOT string value in the database
func (db *SQLiteDB) GetBotStringValue(key string, locale string) (string, error) {
	if locale == "" {
		locale = consts.DefaultLocale
	}
	res := ""
	err := db.QueryRow("SELECT Value FROM `Strings` WHERE `Key` = ? AND `Locale` = ?",
		key, locale).Scan(&res)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetBotStringValue_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return res, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetBotStringValue_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return res, err
	default:
		return res, err
	}
}

//SetBotStringValue sets a value in the bot settings table
func (db *SQLiteDB) SetBotStringValue(key string, value string, locale string) error {
	query, err := db.Exec(
		"INSERT INTO Settings (`Key`, `Value`) VALUES (?,?) "+
			"ON CONFLICT(`Key`, `Locale`) DO UPDATE SET `Value` = Excluded.Value",
		key, value)
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
