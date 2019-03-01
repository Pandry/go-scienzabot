package database

import (
	"database/sql"
	"errors"
	"scienzabot/consts"
	"time"
)

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The usersHelpers.go file focuses on the Users table in the database, that contains
//	all the users subrscribed to the bot

//TODO: Remove user

//AddUser takes a a database.User struct as parameter and insert it in the database
//Only ID, Nickname, Status and Locale will be considered since other ones are supposed to be setted later
func (db *SQLiteDB) AddUser(usr User) error {
	if usr.Locale == "" {
		usr.Locale = consts.DefaultLocale
	}
	var nick sql.NullString
	nick.String = usr.Nickname
	nick.Valid = usr.Nickname != ""

	query, err := db.Exec("INSERT INTO Users (`ID`, `Nickname`, `Status`, `Permissions`, `Locale`) VALUES (?,?,?,?,?)",
		usr.ID, nick, usr.Status, usr.Permissions, usr.Locale)
	if err != nil {
		db.AddLogEvent(Log{Event: "AddUser_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "AddUser_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "AddUser_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GetUser returns a user using the database.user struct
func (db *SQLiteDB) GetUser(userID int) (User, error) {
	var (
		usr  User
		bio  sql.NullString
		nick sql.NullString
	)
	err := db.QueryRow("SELECT `ID`, `Nickname`, `Biography`, `Status`, `LastSeen`, `RegisterDate`, `Permissions`, `Locale` "+
		"FROM Users WHERE `ID`=?", userID).Scan(
		&usr.ID, &nick, &bio, &usr.Status, &usr.LastSeen, &usr.RegisterDate, &usr.Permissions, &usr.Locale)

	usr.Biography = bio.String
	usr.Nickname = nick.String
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetUser_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return usr, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetUser_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return usr, err
	default:
		return usr, err
	}
}

//UpdateUser updates a user data, using a reference the ID.
//All the fields will be used, so make sure that every field of the user struct contains something!
func (db *SQLiteDB) UpdateUser(user User) error {
	if user.Locale == "" {
		user.Locale = consts.DefaultLocale
	}
	query, err := db.Exec("UPDATE Users SET `Nickname` = ?,`Biography` = ?, `Status` = ?, `LastSeen` = ?, `Permissions` = ?, `Locale` = ?  WHERE `ID`=?",
		user.Nickname, user.Biography, user.Status, user.LastSeen, user.Permissions, user.Locale, user.ID)
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateUser_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateUser_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "UpdateUser_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err

}

//UserExists returs a bool that indicates if the user exists or not
func (db *SQLiteDB) UserExists(userID int) bool {
	var dummyval int64
	err := db.QueryRow("SELECT 1 FROM `Users` WHERE `ID`=?", userID).Scan(&dummyval)
	switch {
	case err == sql.ErrNoRows:
		//db.AddLogEvent(Log{Event: "_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return false
	case err != nil:
		db.AddLogEvent(Log{Event: "UserExists_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return false
	default:
		return true
	}

}

//GetUserIDByNickname returns the ID of a user given its nickname.
//Returns sql.ErrNoRows if the user is not present.
//it is safe to just check if the error is nil to confirm if the nick exists
func (db *SQLiteDB) GetUserIDByNickname(nickname string) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT `ID` FROM `Users` WHERE `Nickname` = ?", "a").Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetUserIDByNickname_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return id, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetUserIDByNickname_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return id, err
	default:
		return id, nil
	}
}

//SetUserBiography sets the biography of a user
func (db *SQLiteDB) SetUserBiography(userID int, bio string) error {
	query, err := db.Exec("UPDATE Users SET `Biography` = ? WHERE `ID` = ?", bio, userID)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetUserBiography_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetUserBiography_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetUserBiography_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GetUserBiography returns the biography of a user
func (db *SQLiteDB) GetUserBiography(userID int) (string, error) {
	var bio sql.NullString
	err := db.QueryRow("SELECT `Biography` FROM Users WHERE `ID` = ?", userID).Scan(&bio)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetUserBiography_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return bio.String, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetUserBiography_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return bio.String, err
	default:
		return bio.String, err
	}
}

//SetUserPermissions sets the permissions of a user
func (db *SQLiteDB) SetUserPermissions(userID int, perm int64) error {
	query, err := db.Exec("UPDATE Users SET `Permissions` = ? WHERE `ID` = ?", perm, userID)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetUserPermissions_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetUserPermissions_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetUserPermissions_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GetUserPermissions returns the permissions of a user
func (db *SQLiteDB) GetUserPermissions(userID int) (int64, error) {
	var perm int64
	err := db.QueryRow("SELECT `Permissions` FROM Users WHERE `ID` = ?", userID).Scan(&perm)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetUserPermissions_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return perm, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetUserPermissions_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return perm, err
	default:
		return perm, err
	}
}

//SetUserNickname sets the nickname of a user
func (db *SQLiteDB) SetUserNickname(userID int, userNickname string) error {
	var nick sql.NullString
	nick.String = userNickname
	nick.Valid = userNickname != ""
	query, err := db.Exec("UPDATE Users SET `Nickname` = ? WHERE `ID` = ?", nick, userID)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetUserNickname_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetUserNickname_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetUserNickname_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//SetUserLocale sets the locale of a user
func (db *SQLiteDB) SetUserLocale(userID int, userLocale string) error {
	if userLocale == "" {
		return nil
	}
	query, err := db.Exec("UPDATE Users SET `Locale` = ? WHERE `ID` = ?", userLocale, userID)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetUserLocale_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetUserLocale_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetUserLocale_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GetUserLocale returns the locale of a user or the default one if not found (returning an error)
//It always return a locale
func (db *SQLiteDB) GetUserLocale(userID int) (string, error) {
	locale := consts.DefaultLocale
	err := db.QueryRow("SELECT `Locale` FROM Users WHERE `ID` = ?", userID).Scan(&locale)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetUserLocale_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return locale, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetUserLocale_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return locale, err
	default:
		return locale, err
	}
}

//UpdateUserLastSeen updates the lastseen field
func (db *SQLiteDB) UpdateUserLastSeen(userID int, lastSeen time.Time) error {
	//lastSeenString := lastSeen.Format("dd/MM/YYYY hh:mm:ss")
	lastSeenString := lastSeen.Format(consts.TimeFormatString)
	query, err := db.Exec("UPDATE Users SET `LastSeen` = ? WHERE `ID` = ?", lastSeenString, userID)
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateUserLastSeen_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateUserLastSeen_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "UpdateUserLastSeen_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//UpdateUserLastSeenToNow updates the lastseen field to the actual time
func (db *SQLiteDB) UpdateUserLastSeenToNow(userID int) error {
	return db.UpdateUserLastSeen(userID, time.Now())
}
