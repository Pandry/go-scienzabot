package database

import (
	"database/sql"
	"errors"
	"time"
)

/*ID'  INTEGER NOT NULL PRIMARY KEY,
'Nickname'  TEXT UNIQUE,
'Biography'  TEXT,
'Status'  INTEGER DEFAULT 0,
'LastSeen*/

//TODO: Remove user

//AddUser takesa a database.User struct as parameter and insert it in the database
//Only ID, Nickname and Status will be considered since other ones are supposed to be setted later
func (db *SQLiteDB) AddUser(usr User) error {
	stmt, err := db.Prepare("INSERT INTO Users (`ID`, `Nickname`, `Status`) VALUES (?,?,?)")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "AddUser_QueryFailed", Message: "Impossible to create the AddUser preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(usr.ID, usr.Nickname, usr.Status)
	if err != nil {
		db.AddLogEvent(Log{Event: "AddUser_ExecutionQueryFailed", Message: "Impossible to execute the AddUser query", Error: err.Error()})
	}
	return err
}

//GetUser returns a user using the database.user struct
func (db *SQLiteDB) GetUser(userID int64) (User, error) {
	result := User{}
	result.ID = userID
	//We're prepaing a query to execute later
	/*
		'LastSeen*/
	//stmt, err := db.Prepare("SELECT `Nickname`,`Biography`,`Status`,`LastSeen` FROM `Users` WHERE `ID` = ?")
	stmt, err := db.Prepare("SELECT `Nickname`,`Biography`,`Status` FROM `Users` WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetUser_QueryFailed", Message: "Impossible to create the GetUser preparation query", Error: err.Error()})
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
		db.AddLogEvent(Log{Event: "RequestedUserID_NotFound", Message: "Requested a user not found in the database", Error: err.Error()})
		return result, err

	case err != nil:
		db.AddLogEvent(Log{Event: "RequestedUserID_NotFoundUnknown", Message: "Requested a user not found in the database but the error is unknown", Error: err.Error()})
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
			result.Status = status.Int64
		}
		return result, err
	}
}

//UpdateUser updates a user data, using a reference the ID.
//All the fields will be used, so make sure that every field of the user struct contains something!
func (db *SQLiteDB) UpdateUser(user User) error {
	stmt, err := db.Prepare("UPDATE Users SET `Nickname` = ?,`Biography` = ?, `Status` = ?, `LastSeen` = ?  WHERE `ID`=?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UpdateUser_QueryFailed", Message: "Impossible to create the UpdateUser preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	res, err := stmt.Exec(user.Nickname, user.Biography, user.Status, user.LastSeen, user.ID)
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateUser_NotFoundUnknown", Message: "Requested a user edit, but the execution of the query returned an error", Error: err.Error()})
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateUser_NotFoundUnknown", Message: "Requested a user edit, but the rows afftected method returned an error", Error: err.Error()})
		return err
	}
	if rows > 0 {
		return nil
	}
	db.AddLogEvent(Log{Event: "UpdateUser_NoRowsAffected", Message: "Requested a user edit, but the ID was not found.", Error: err.Error()})
	return errors.New("UserDoesNotExistError")

}

//UserExists returs a bool that indicates if the user exists or not
func (db *SQLiteDB) UserExists(userID int) bool {
	stmt, err := db.Prepare("SELECT 1 FROM `Users` WHERE `ID`=?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UserExists_QueryFailed", Message: "Impossible to create the UserExists preparation query", Error: err.Error()})
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
		db.AddLogEvent(Log{Event: "UserExists_UnknownError", Message: "Requested a user ID not present in the database but the error is unknown", Error: err.Error()})
		return false

	default:
		//Success
		if res.Valid && res.Int64 == 1 {
			return true
		}
		return false
	}
}

//GetUserIDByNickname returns the ID of a user given its nickname.
//Returns sql.ErrNoRows if the user is not present.
//it is safe to just check if the error is nil to confirm if the nick exists
func (db *SQLiteDB) GetUserIDByNickname(nickname string) (int, error) {
	stmt, err := db.Prepare("SELECT `ID` FROM `Users` WHERE `Nickname` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetUserIDByNickname_QueryFailed", Message: "Impossible to create the GetUserIDByNickname preparation query", Error: err.Error()})
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
		db.AddLogEvent(Log{Event: "RequestedNickname_DontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return 0, err

	case err != nil:
		db.AddLogEvent(Log{Event: "RequestedNickname_DontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return 0, err

	default:
		//result.LastSeen, err = time.Parse("2006-01-02 20:50:59", lastseen)
		//Success
		return id, err
	}
}

//SetUserBiography sets the biography of a user
func (db *SQLiteDB) SetUserBiography(userID int, bio string) error {
	stmt, err := db.Prepare("UPDATE Users SET `Biography` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "SetUserBiography_QueryFailed", Message: "Impossible to create the SetUserBiography preparation query", Error: err.Error()})
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
		db.AddLogEvent(Log{Event: "SetUserBiography_UserDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "SetUserBiography_UserDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}
}

//GetUserBiography returns the biography of a user
func (db *SQLiteDB) GetUserBiography(userID int) (string, error) {
	stmt, err := db.Prepare("SELECT `Biography` FROM Users WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetUserBiography_QueryFailed", Message: "Impossible to create the GetUserBiography preparation query", Error: err.Error()})
		return "", err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	res := stmt.QueryRow(userID)

	var bio sql.NullString
	err = res.Scan(&bio)

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "GetUserBiography_UserDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return "", err

	case err != nil:
		db.AddLogEvent(Log{Event: "GetUserBiography_UserDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return "", err
	}

	if !bio.Valid {
		//Bio is null
		return "", nil
	}
	return bio.String, nil

}

//SetUserNickname sets the nickname of a user
func (db *SQLiteDB) SetUserNickname(userID int, userNickname string) error {
	stmt, err := db.Prepare("UPDATE Users SET `Nickname` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "SetUserNickname_QueryFailed", Message: "Impossible to create the SetUserBiography preparation query", Error: err.Error()})
		return err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	_, err = stmt.Exec(userNickname, userID)

	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "SetUserNickname_UserDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "SetUserNickname_UserDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}

}

//UpdateUserLastSeen updates the lastseen field
func (db *SQLiteDB) UpdateUserLastSeen(userID int, lastSeen time.Time) error {
	lastSeenString := lastSeen.Unix()
	stmt, err := db.Prepare("UPDATE Users SET `LastSeen` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UpdateUserLastSeen_QueryFailed", Message: "Impossible to create the UpdateUserLastSeen preparation query", Error: err.Error()})
		return err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	_, err = stmt.Exec(lastSeenString, userID)

	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "UpdateUserLastSeen_UserDontExists", Message: "Requested an ID not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "UpdateUserLastSeen_UserDontExistsUnknown", Message: "Requested an ID not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}
}

//UpdateUserLastSeenToNow updates the lastseen field to the actual time
func (db *SQLiteDB) UpdateUserLastSeenToNow(userID int) error {
	return db.UpdateUserLastSeen(userID, time.Now())
}
