package database

import "errors"

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
func (db *SQLiteDB) SetPermissions(userID int, groupID int, permissions int) error {
	stmt, err := db.Prepare("INSERT INTO Permissions (`User`, `Group` , `Permission`) VALUES (?,?,?) ON CONFLICT(`User`,`Group`,`Permission`) DO UPDATE SET `Permission` = Excluded.Permission")
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
	return errors.New("No changes to the database was made")
}
