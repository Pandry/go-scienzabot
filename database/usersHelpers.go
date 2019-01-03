package database

import (
	"database/sql"
)

/*ID'  INTEGER NOT NULL PRIMARY KEY,
'Nickname'  TEXT UNIQUE,
'Biography'  TEXT,
'Status'  INTEGER DEFAULT 0,
'LastSeen*/

//AddUser takesa a database.User struct as parameter and insert it in the database
//Only ID, Nickname and Status will be considered since other ones are supposed to be setted later
func (db *SQLiteDB) AddUser(usr User) error {
	stmt, err := db.Prepare("INSERT INTO Users (`ID`, `Nickname`, `Status`) VALUES (?,?,?)")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "AddUserQueryFailed", Message: "Impossible to create the AddUser preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	stmt.Exec(usr.ID, usr.Nickname, usr.Status)

	return nil
}

//GetUserIDByNickname returns the ID of a user given its nickname.
//Returns sql.ErrNoRows if the user is not present.
//it is safe to just check if the error is nil to confirm if the nick exists
func (db *SQLiteDB) GetUserIDByNickname(nickname string) (int, error) {
	stmt, err := db.Prepare("SELECT `ID` FROM `Users` WHERE `Nickname` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetUserIDByNicknameQueryFailed", Message: "Impossible to create the GetUserIDByNickname preparation query", Error: err.Error()})
		return 0, err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	query := stmt.QueryRow(nickname)

	var id int
	err = query.Scan(&id)
	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "RequestedNicknameDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return 0, err

	case err != nil:
		db.AddLogEvent(Log{Event: "RequestedNicknameDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return 0, err

	default:
		//result.LastSeen, err = time.Parse("2006-01-02 20:50:59", lastseen)
		//Success
		return id, err
	}
}

//GetUser returns a user using the database.user struct
func (db *SQLiteDB) GetUser(userID int) (User, error) {
	result := User{}
	result.ID = userID
	//We're prepaing a query to execute later
	/*
		'LastSeen*/
	//stmt, err := db.Prepare("SELECT `Nickname`,`Biography`,`Status`,`LastSeen` FROM `Users` WHERE `ID` = ?")
	stmt, err := db.Prepare("SELECT `Nickname`,`Biography`,`Status` FROM `Users` WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetUserQueryFailed", Message: "Impossible to create the GetUser preparation query", Error: err.Error()})
		return result, err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	query := stmt.QueryRow(userID)

	//We than create some types that are nullable (some filends in the DB can be null and the default types
	//Does not support null values)
	var nickname, biography sql.NullString
	var status sql.NullInt64

	//We then scan the query
	err = query.Scan(&nickname, &biography, &status)
	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "RequestedUserIDNotFound", Message: "Requested a user not found in the database", Error: err.Error()})
		return result, err

	case err != nil:
		db.AddLogEvent(Log{Event: "RequestedUserIDNotFoundUnknown", Message: "Requested a user not found in the database but the error is unknown", Error: err.Error()})
		return result, err

	default:
		//result.LastSeen, err = time.Parse("2006-01-02 20:50:59", lastseen)
		//Success
		if nickname.Valid {
			result.Nickname = nickname.String
		}
		if biography.Valid {
			result.Biography = biography.String
		}
		if status.Valid {
			result.Status = int(status.Int64)
		}
		return result, err
	}
}

//SetUserBiography sets the biography of a user
func (db *SQLiteDB) SetUserBiography(userID int, bio string) error {
	stmt, err := db.Prepare("UPDATE Users SET `Biography` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "SetUserBiographyQueryFailed", Message: "Impossible to create the SetUserBiography preparation query", Error: err.Error()})
		return err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	_, err = stmt.Exec(bio, userID)

	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "SetUserBiographyUserDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "SetUserBiographyUserDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}
}

//UserExists returs a bool that indicates if the user exists or not
func (db *SQLiteDB) UserExists(userID int) bool {
	stmt, err := db.Prepare("SELECT 1 FROM `Users` WHERE `ID`=?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UserExistsQueryFailed", Message: "Impossible to create the UserExists preparation query", Error: err.Error()})
		return false
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	qry := stmt.QueryRow(userID)

	var res sql.NullInt64
	err = qry.Scan(&res)
	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		return false

	case err != nil:
		db.AddLogEvent(Log{Event: "UserExistsUnknownError", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return false

	default:
		//Success
		if res.Valid && res.Int64 == 1 {
			return true
		}
		return false
	}
}
