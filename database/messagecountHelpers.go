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

// The messagecountHelpers.go file focuses on the MessageCount table in the database.
// The messageCount feaure is used to count the number of messages of each subscribed user
//	in every group the bot is in

//GetMessageCount returns the message number of a user in a group
func (db *SQLiteDB) GetMessageCount(user int64, group int64) (int64, error) {
	var messageCount int64
	err := db.QueryRow("SELECT MessageCount FROM MessageCount WHERE `UserID` AND `GroupID`", user, group).Scan(&messageCount)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetMessageCount_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return messageCount, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetMessageCount_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return messageCount, err
	default:
		return messageCount, nil
	}
}

//GetListsInvokedCount returns the number of lists invoked by a user in a group
func (db *SQLiteDB) GetListsInvokedCount(user int, group int64) (int64, error) {
	var listInvocations int64
	err := db.QueryRow("SELECT ListsInvoked FROM MessageCount WHERE `UserID` = ? AND `GroupID` = ?", user, group).Scan(&listInvocations)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetListsInvokedCount_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return listInvocations, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetListsInvokedCount_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return listInvocations, err
	default:
		return listInvocations, nil
	}
}

//SetMessageCount sets the message of a user in a group
func (db *SQLiteDB) SetMessageCount(user int64, group int64, messageCount int64) error {
	query, err := db.Exec(
		"INSERT INTO MessageCount (`UserID`, `GroupID`, `MessageCount`) VALUES (?,?,?) "+
			"ON CONFLICT(`UserID`, `GroupID`) DO UPDATE SET `MessageCount` = Excluded.MessageCount",
		user, group, messageCount)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetMessageCount_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetMessageCount_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetMessageCount_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//SetListsInvokedCount sets the number of lists invoked of a user
func (db *SQLiteDB) SetListsInvokedCount(user int, group int64, listsInvoked int64) error {
	query, err := db.Exec(
		"INSERT INTO MessageCount (`UserID`, `GroupID`, `ListsInvoked`) VALUES (?,?,?) "+
			"ON CONFLICT(`UserID`, `GroupID`) DO UPDATE SET `ListsInvoked` = Excluded.ListsInvoked",
		user, group, listsInvoked)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetListsInvokedCount_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetListsInvokedCount_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetListsInvokedCount_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//IncrementMessageCount increments by 1 the number of messages from a user
func (db *SQLiteDB) IncrementMessageCount(user int64, group int64) error {
	msgCnt, err := db.GetMessageCount(user, group)
	/*if err != nil {
		//User may not exist yet in DB
		db.AddLogEvent(Log{Event: "IncrementMessageCount_CannotGetMessageCount", Message: "user may not exist", Error: err.Error()})
	}*/
	err = db.SetMessageCount(user, group, msgCnt+1)
	return err
}

//IncrementListsInvokedCount increments by 1 the number of messages from a user
func (db *SQLiteDB) IncrementListsInvokedCount(user int, group int64) error {
	lstsCnt, err := db.GetListsInvokedCount(user, group)
	err = db.SetListsInvokedCount(user, group, lstsCnt+1)
	return err
}

//GetUserGroups returns the groups
func (db *SQLiteDB) GetUserGroups(user int) ([]Group, error) {
	gprs := make([]Group, 0)
	rows, err := db.Query("SELECT Groups.ID, Groups.Title,Groups.Status, Groups.Locale, Groups.Ref FROM MessageCount INNER JOIN Groups ON MessageCount.GroupID = Groups.ID  WHERE `UserID`=?", user)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetUserGroups_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return gprs, err
	}
	for rows.Next() {
		var grp Group
		if err = rows.Scan(&grp.ID, &grp.Title, &grp.Status, &grp.Locale, &grp.Ref); err != nil {
			db.AddLogEvent(Log{Event: "GetUserGroups_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			gprs = append(gprs, grp)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetUserGroups_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetUserGroups_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}
	return gprs, err
}

//UpdateLastInvocation updates the lastseen field
func (db *SQLiteDB) UpdateLastInvocation(userID int, groupID int64, lstInv time.Time) error {
	lastInvocation := lstInv.Format(consts.TimeFormatString)
	query, err := db.Exec("UPDATE MessageCount SET `LatestListInvocation` = ? WHERE `UserID` = ? AND `GroupID` = ?", lastInvocation, userID, groupID)
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateLastInvocation_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateLastInvocation_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "UpdateLastInvocation_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GetLastListInvocation returns the last time when the user invoked a list
func (db *SQLiteDB) GetLastListInvocation(user int, group int64) (time.Time, error) {
	var listInvocation time.Time
	var timeStr sql.NullString
	err := db.QueryRow("SELECT LatestListInvocation FROM MessageCount WHERE `UserID` AND `GroupID`", user, group).Scan(&timeStr)
	listInvocation, _ = time.Parse(consts.TimeFormatString, timeStr.String)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetLastListInvocation_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return listInvocation, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetLastListInvocation_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return listInvocation, err
	default:
		return listInvocation, nil
	}
}
