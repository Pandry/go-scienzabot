package database

import (
	"database/sql"
	"errors"
	"scienzabot/consts"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS 'Users' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Nickname'  TEXT UNIQUE,
	'Biography'  TEXT,
	'Status'  INTEGER NOT NULL DEFAULT 0,
	'LastSeen'  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	'RegisterDate' TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

*/

//TODO: Remove user

//AddUser takes a a database.User struct as parameter and insert it in the database
//Only ID, Nickname, Status and Locale will be considered since other ones are supposed to be setted later
func (db *SQLiteDB) AddUser(usr User) error {
	if usr.Locale == "" {
		usr.Locale = consts.DefaultLocale
	}
	var nick sql.NullString
	nick.String = usr.Nickname
	nick.Valid = usr.Nickname == ""

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
		usr User
		bio sql.NullString
	)
	err := db.QueryRow("SELECT `ID`, `Nickname`, `Biography`, `Status`, `LastSeen`, `RegisterDate`, `Permissions`, `Locale` "+
		"FROM Users WHERE `ID`=?", userID).Scan(
		&usr.ID, &usr.Nickname, &bio, &usr.Status, &usr.LastSeen, &usr.RegisterDate, &usr.Permissions, &usr.Locale)

	if bio.Valid {
		usr.Biography = bio.String
	}
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
	nick.Valid = userNickname == ""
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
	var loc sql.NullString
	loc.String = userLocale
	loc.Valid = userLocale == ""
	query, err := db.Exec("UPDATE Users SET `Locale` = ? WHERE `ID` = ?", loc, userID)
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

//UpdateUserLastSeen updates the lastseen field
func (db *SQLiteDB) UpdateUserLastSeen(userID int, lastSeen time.Time) error {
	lastSeenString := lastSeen.Format("dd/MM/YYYY hh:mm:ss")
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
