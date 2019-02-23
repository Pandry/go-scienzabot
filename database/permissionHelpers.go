package database

import (
	"database/sql"
	"errors"
)

/*
CREATE TABLE IF NOT EXISTS 'Permissions' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'User'	INTEGER NOT NULL,
	'Group'	INTEGER,
	'Permission' INTEGER DEFAULT 0,
	FOREIGN KEY('User') REFERENCES Users('ID'),
	FOREIGN KEY('Group') REFERENCES Groups('ID'),
	CONSTRAINT con_perm_user_group_perm_unique UNIQUE ('User','Group','Permission')
);
*/

//SetPermissions sets the permissions of  user in a group
func (db *SQLiteDB) SetPermissions(prm Permission) error {

	query, err := db.Exec("INSERT INTO Permissions (`UserID`, `GroupID` , `Permission`) VALUES (?,?,?) "+
		"ON CONFLICT(`UserID`,`GroupID`,`Permission`) DO UPDATE "+
		"SET `Permission` = Excluded.Permission", prm.UserID, prm.GroupID, prm.Permission)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetPermissions_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetPermissions_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetPermissions_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err

	/*stmt, err := db.Prepare("INSERT INTO Permissions (`User`, `Group` , `Permission`) VALUES (?,?,?) ON CONFLICT(`User`,`Group`,`Permission`) DO UPDATE SET `Permission` = Excluded.Permission")
	if err != nil {
		db.AddLogEvent(Log{Event: "SetPermissions_QueryFailed", Message: "The query for the SetPermissions function failed", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	rows, err := stmt.Exec(userID, groupID, permissions)

	if err != nil {
		db.AddLogEvent(Log{Event: "SetPermissions_NotFoundUnknown", Message: "The execution of the query for the SetPermissions function failed", Error: err.Error()})
		return err
	}

	res, err := rows.RowsAffected()

	if err != nil {
		db.AddLogEvent(Log{Event: "SetPermissions_ExecutionQueryError", Message: "The fetching of the query results for the SetPermissions function failed", Error: err.Error()})
		return err
	}
	if res > 0 {
		return nil
	}
	db.AddLogEvent(Log{Event: "SetPermissions_NotChangesMade", Message: "No changes was made to the database!", Error: err.Error()})
	return errors.New("No changes to the database was made")*/
}

//RemoveAllGroupPermissions removes all the admins in a group
func (db *SQLiteDB) RemoveAllGroupPermissions(groupID int64) error {

	query, err := db.Exec("DELETE FROM Permissions WHERE `GroupID` =?", groupID)
	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveAllGroupPermissions_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveAllGroupPermissions_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "RemoveAllGroupPermissions_Info_NoRowsAffected", Message: "No rows affected"})
		return nil
		// ACTUALLY, if the group has no permissions yet, it's ok to have0 results (but we still log 'em)
		//return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GetPermission returns the permission of a user given its id and the group
func (db *SQLiteDB) GetPermission(userID int64, groupID int64) (int, error) {
	var perm int64
	err := db.QueryRow("SELECT `Permission` FROM Permissions WHERE `ID`=?", "a").Scan(&perm)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetPermission_Info_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return int(perm), nil
	case err != nil:
		db.AddLogEvent(Log{Event: "GetPermission_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return int(perm), nil
	default:
		return int(perm), nil
	}
}

//GetPrivilegedUsers returns an array of database.Permission type ofall the users in the database
func (db *SQLiteDB) GetPrivilegedUsers(groupID int64) ([]Permission, error) {
	rows, err := db.Query("SELECT `Permission` FROM Permissions WHERE `GroupID`=?", groupID)
	defer rows.Close()
	prms := make([]Permission, 0)
	if err != nil {
		db.AddLogEvent(Log{Event: "GetPrivilegedUsers_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return prms, err
	}
	for rows.Next() {
		var (
			id, userID, groupID, perm int64
		)
		if err = rows.Scan(&id, &userID, &groupID, &perm); err != nil {
			db.AddLogEvent(Log{Event: "GetPrivilegedUsers_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			prms = append(prms, Permission{ID: id, UserID: userID, GroupID: groupID, Permission: perm})
		}
	}
	if err == sql.ErrNoRows {
		return prms, nil
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetPrivilegedUsers_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetPrivilegedUsers_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}
	return prms, err
}
