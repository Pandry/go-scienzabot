package database

import "database/sql"

/*
type Bookmark struct {
	ID             int64
	UserID         int64
	GroupID        int64
	MessageID      int64
	Alias          string
	Status         int64
	MessageContent string
}
*/

//CreateBookmark takes a database.Bookmark struct as parameter and insert it in the database
//The ID will not be considered, since it's automatically inrted in database. ALl the other values will be inserted
func (db *SQLiteDB) CreateBookmark(bkm Bookmark) error {
	stmt, err := db.Prepare("INSERT INTO Bookmarks (`UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent`) VALUES (?,?,?,?,?,?)")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "CreateBookmark_QueryFailed", Message: "Impossible to create the CreateBookmark preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(bkm.UserID, bkm.UserID, bkm.GroupID, bkm.MessageID, bkm.Alias, bkm.Status, bkm.MessageContent)
	if err != nil {
		db.AddLogEvent(Log{Event: "CreateBookmark_ExecutionQueryFailed", Message: "Impossible to execute the CreateBookmark query", Error: err.Error()})
	}
	return err
}

//DeleteBookmark takes a Bookmark ID and deletes it
func (db *SQLiteDB) DeleteBookmark(bkmID int) error {
	stmt, err := db.Prepare("DELETE FROM Bookmarks WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "DeleteBookmark_QueryFailed", Message: "Impossible to create the DeleteBookmark preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(bkmID)
	if err != nil {
		db.AddLogEvent(Log{Event: "DeleteBookmark_ExecutionQueryFailed", Message: "Impossible to execute the DeleteBookmark query", Error: err.Error()})
	}
	return err
}

//RenameBookmark takes a Bookmark ID and a string and renames  a bookmark
func (db *SQLiteDB) RenameBookmark(bkmID int, newAlias string) error {
	stmt, err := db.Prepare("UPDATE Bookmarks SET `Alias`=? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "RenameBookmark_QueryFailed", Message: "Impossible to create the RenameBookmark preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(newAlias, bkmID)
	if err != nil {
		db.AddLogEvent(Log{Event: "RenameBookmark_ExecutionQueryFailed", Message: "Impossible to execute the RenameBookmark query", Error: err.Error()})
	}
	return err
}

//GetAllBookmarks returns all the bookmark in the database
func (db *SQLiteDB) GetAllBookmarks() ([]Bookmark, error) {
	stmt, err := db.Prepare("SELECT `ID`, `UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent` FROM Bookmarks")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "RenameBookmark_QueryFailed", Message: "Impossible to create the RenameBookmark preparation query", Error: err.Error()})
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "RequestedNickname_DontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return nil, err

	case err != nil:
		db.AddLogEvent(Log{Event: "RequestedNickname_DontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return nil, err
	}

	var bkms []Bookmark

	var id, userID, groupID, messageID, status int64
	var messageContent, alias string
	for rows.Next() {
		err = rows.Scan(&id, &userID, &groupID, &messageID, &alias, &status, &messageContent)
		//And check for errors
		if err != nil {
			db.AddLogEvent(Log{Event: "GetAllBookmarks_QueryScanError", Message: "An error verified while scanning a result row", Error: err.Error()})
		} else {
			bkms = append(bkms, Bookmark{ID: id, UserID: userID, GroupID: groupID, MessageID: messageID, Alias: alias, Status: status, MessageContent: messageContent})
		}
	}
	return bkms, err
}

//GetUserBookmarks returns all the bookmarks of a user
func (db *SQLiteDB) GetUserBookmarks(iUserID int64) ([]Bookmark, error) {
	stmt, err := db.Prepare("SELECT `ID`, `UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent` FROM Bookmarks WHERE `UserID`=?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetUserBookmarks_QueryFailed", Message: "Impossible to create the GetUserBookmarks preparation query", Error: err.Error()})
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(iUserID)

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "GetUserBookmarks_DontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return nil, err

	case err != nil:
		db.AddLogEvent(Log{Event: "GetUserBookmarks_DontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return nil, err
	}

	var bkms []Bookmark

	var id, userID, groupID, messageID, status int64
	var messageContent, alias string
	for rows.Next() {
		err = rows.Scan(&id, &userID, &groupID, &messageID, &alias, &status, &messageContent)
		//And check for errors
		if err != nil {
			db.AddLogEvent(Log{Event: "GetUserBookmarks_QueryScanError", Message: "An error verified while scanning a result row", Error: err.Error()})
		} else {
			bkms = append(bkms, Bookmark{ID: id, UserID: userID, GroupID: groupID, MessageID: messageID, Alias: alias, Status: status, MessageContent: messageContent})
		}
	}
	return bkms, err
}

//GetUserGroupBookmarks returns the bookmarks of a user in a given
func (db *SQLiteDB) GetUserGroupBookmarks(iUserID int64, iGroupID int64) ([]Bookmark, error) {
	stmt, err := db.Prepare("SELECT `ID`, `UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent` FROM Bookmarks WHERE `UserID`=? AND `GroupID`=?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetUserGroupBookmarks_QueryFailed", Message: "Impossible to create the GetUserGroupBookmarks preparation query", Error: err.Error()})
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(iUserID, iGroupID)

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "GetUserGroupBookmarks_DontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return nil, err

	case err != nil:
		db.AddLogEvent(Log{Event: "GetUserGroupBookmarks_DontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return nil, err
	}

	var bkms []Bookmark

	var id, userID, groupID, messageID, status int64
	var messageContent, alias string
	for rows.Next() {
		err = rows.Scan(&id, &userID, &groupID, &messageID, &alias, &status, &messageContent)
		//And check for errors
		if err != nil {
			db.AddLogEvent(Log{Event: "GetUserGroupBookmarks_QueryScanError", Message: "An error verified while scanning a result row", Error: err.Error()})
		} else {
			bkms = append(bkms, Bookmark{ID: id, UserID: userID, GroupID: groupID, MessageID: messageID, Alias: alias, Status: status, MessageContent: messageContent})
		}
	}
	return bkms, err
}
