package database

import (
	"errors"
)

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The bookmarksHelpers.go file focuses on the Bookmarks table in the database.
// The bookmark feature is basically a feature that saves a message to read it later moment.
// It isn't implemented yet.

//CreateBookmark takes a database.Bookmark struct as parameter and insert it in the database
//The ID will not be considered, since it's automatically inrted in database. ALl the other values will be inserted
func (db *SQLiteDB) CreateBookmark(bkm Bookmark) error {
	query, err := db.Exec("INSERT INTO Bookmarks (`UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent`) VALUES (?,?,?,?,?,?)",
		bkm.UserID, bkm.GroupID, bkm.MessageID, bkm.Alias, bkm.Status, bkm.MessageContent)
	if err != nil {
		db.AddLogEvent(Log{Event: "CreateBookmark_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "CreateBookmark_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "CreateBookmark_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//DeleteBookmark takes a Bookmark ID and deletes it
func (db *SQLiteDB) DeleteBookmark(bkmID int) error {
	query, err := db.Exec("DELETE FROM Bookmarks WHERE `ID` = ?", bkmID)
	if err != nil {
		db.AddLogEvent(Log{Event: "DeleteBookmark_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "DeleteBookmark_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "DeleteBookmark_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//RenameBookmark takes a Bookmark ID and a string and renames  a bookmark
func (db *SQLiteDB) RenameBookmark(bkmID int, newAlias string) error {
	query, err := db.Exec("UPDATE Bookmarks SET `Alias`=? WHERE `ID` = ?", newAlias, bkmID)
	if err != nil {
		db.AddLogEvent(Log{Event: "RenameBookmark_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "RenameBookmark_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "RenameBookmark_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GetAllBookmarks returns all the bookmark in the database
func (db *SQLiteDB) GetAllBookmarks() ([]Bookmark, error) {
	rows, err := db.Query("SELECT `ID`, `UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent` FROM Bookmarks")
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetAllBookmarks_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	bkms := make([]Bookmark, 0)
	for rows.Next() {
		var (
			id, userID, groupID, messageID, status int64
			messageContent, alias                  string
		)
		if err = rows.Scan(&id, &userID, &groupID, &messageID, &alias, &status, &messageContent); err != nil {
			db.AddLogEvent(Log{Event: "GetAllBookmarks_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			bkms = append(bkms, Bookmark{ID: id, UserID: userID, GroupID: groupID, MessageID: messageID, Alias: alias, Status: status, MessageContent: messageContent})
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetAllBookmarks_RowNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetAllBookmarks_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return bkms, err
}

//GetUserBookmarks returns all the bookmarks of a user ordered by groupID
func (db *SQLiteDB) GetUserBookmarks(userid int) ([]Bookmark, error) {
	iUserID := int64(userid)
	rows, err := db.Query("SELECT `ID`, `UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent` FROM Bookmarks WHERE `UserID`=? ORDER BY `GroupID`", iUserID)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetUserBookmarks_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	bkms := make([]Bookmark, 0)
	for rows.Next() {
		var (
			id, userID, groupID, messageID, status int64
			messageContent, alias                  string
		)
		if err = rows.Scan(&id, &userID, &groupID, &messageID, &alias, &status, &messageContent); err != nil {
			db.AddLogEvent(Log{Event: "GetUserBookmarks_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			bkms = append(bkms, Bookmark{ID: id, UserID: userID, GroupID: groupID, MessageID: messageID, Alias: alias, Status: status, MessageContent: messageContent})
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetUserBookmarks_RowNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetUserBookmarks_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return bkms, err

}

//GetUserGroupBookmarks returns the bookmarks of a user in a given
func (db *SQLiteDB) GetUserGroupBookmarks(iUserID int, iGroupID int64) ([]Bookmark, error) {

	rows, err := db.Query("SELECT `ID`, `UserID`, `GroupID`, `MessageID`, `Alias`, `Status`, `MessageContent` FROM Bookmarks WHERE `UserID`=? AND `GroupID`=?", iUserID, iGroupID)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetUserGroupBookmarks_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	bkms := make([]Bookmark, 0)
	for rows.Next() {
		var (
			id, userID, groupID, messageID, status int64
			messageContent, alias                  string
		)
		if err = rows.Scan(&id, &userID, &groupID, &messageID, &alias, &status, &messageContent); err != nil {
			db.AddLogEvent(Log{Event: "GetUserGroupBookmarks_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			bkms = append(bkms, Bookmark{ID: id, UserID: userID, GroupID: groupID, MessageID: messageID, Alias: alias, Status: status, MessageContent: messageContent})
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetUserGroupBookmarks_RowNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetUserGroupBookmarks_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return bkms, err
}
