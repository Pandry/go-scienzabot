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

// The subscriptionsHelpers.go file focuses on the Subscriptions table in the database, that is
//	supposed to keep track of all the users subrscribed to which list

//AddSubscription adds a subscription. takes as a parameter the userID and the listID
func (db *SQLiteDB) AddSubscription(userID int, listID int) error {
	query, err := db.Exec("INSERT INTO Subscriptions (`ListID`, `UserID`) VALUES (?,?)",
		listID, userID)
	if err != nil {
		db.AddLogEvent(Log{Event: "AddSubscription_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "AddSubscription_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "AddSubscription_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//RemoveSubscriptionByID deletes a subscription given its ID
func (db *SQLiteDB) RemoveSubscriptionByID(subID int) error {
	query, err := db.Exec("DELETE FROM Subscriptions WHERE ID = ?", subID)
	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveSubscription_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveSubscription_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "RemoveSubscription_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//RemoveSubscriptionByListAndUserID deletes a subscription given the list and the user's IDs
func (db *SQLiteDB) RemoveSubscriptionByListAndUserID(listID int, userID int) error {
	query, err := db.Exec("DELETE FROM Subscriptions WHERE `ListID` = ? AND `UserID`=?", listID, userID)
	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveSubscriptionByListAndUserID_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveSubscriptionByListAndUserID_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "RemoveSubscriptionByListAndUserID_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GellAllSubscriptions returns a slice containing all the users in a given list
func (db *SQLiteDB) GellAllSubscriptions() ([]Subscription, error) {

	rows, err := db.Query("SELECT `ID`,`ListID`,`UserID` FROM Subscriptions")
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GellAllSubscriptions_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	subs := make([]Subscription, 0)
	for rows.Next() {
		var sb Subscription
		if err = rows.Scan(&sb.ID, &sb.ListID, &sb.UserID); err != nil {
			db.AddLogEvent(Log{Event: "GellAllSubscriptions_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			subs = append(subs, sb)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GellAllSubscriptions_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GellAllSubscriptions_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return subs, err
}

//GetSubscribedUsers returns a slice containing all the users in a given list
func (db *SQLiteDB) GetSubscribedUsers(lstID int64) ([]Subscription, error) {
	rows, err := db.Query("SELECT `ID`,`ListID`,`UserID` FROM Subscriptions WHERE `ListID`=?", lstID)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	subs := make([]Subscription, 0)
	for rows.Next() {
		var sb Subscription
		if err = rows.Scan(&sb.ID, &sb.ListID, &sb.UserID); err != nil {
			db.AddLogEvent(Log{Event: "GetSubscribedUsers_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			subs = append(subs, sb)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return subs, err
}

//GetListSubscribers returns a database.Subscription slice containing all the list a user is subscribed to
func (db *SQLiteDB) GetListSubscribers(usrID int64) ([]Subscription, error) {
	rows, err := db.Query("SELECT `ID`,`ListID`,`UserID` FROM Subscriptions WHERE `UserID`=?", usrID)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetListSubscribers_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	subs := make([]Subscription, 0)
	for rows.Next() {
		var sb Subscription
		if err = rows.Scan(&sb.ID, &sb.ListID, &sb.UserID); err != nil {
			db.AddLogEvent(Log{Event: "GetListSubscribers_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			subs = append(subs, sb)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetListSubscribers_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetListSubscribers_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return subs, err
}

//GetUserGroupListsWithLimits returns a database.Subscription slice containing all the list a user is subscribed to in a specif group
func (db *SQLiteDB) GetUserGroupListsWithLimits(usrID int64, grpID int64, limit int, offset int) ([]List, error) {
	rows, err := db.Query("SELECT Lists.ID, `Name`, `Properties`, `Parent`, `LatestInvocation`, `CreationDate` FROM Lists JOIN Subscriptions ON Lists.ID = Subscriptions.ListID WHERE "+
		"Subscriptions.UserID=? AND Lists.GroupID=? LIMIT ? OFFSET ?", usrID, grpID, limit, offset)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetUserLists_ErrorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	lists := make([]List, 0)
	for rows.Next() {
		var lst List
		var parent sql.NullInt64
		var linvdate sql.NullString
		if err = rows.Scan(&lst.ID, &lst.Name, &lst.Properties, &parent, &linvdate, &lst.CreationDate); err != nil {
			db.AddLogEvent(Log{Event: "GetUserLists_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			lst.Parent = parent.Int64
			lastInvTime, err := time.Parse(consts.TimeFormatString, linvdate.String)
			if err != nil {
				lastInvTime = time.Unix(0, 0)
			}
			lst.LatestInvocation = lastInvTime
			lists = append(lists, lst)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetUserLists_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetUserLists_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return lists, err
}

//GetUserGroupLists returns a database.Subscription slice containing all the list a user is subscribed to in a specif group
func (db *SQLiteDB) GetUserGroupLists(usrID int64, grpID int64) ([]List, error) {

	rows, err := db.Query("SELECT `ID`, `Name`, `Properties`, `Parent`, `LatestInvocation`, `CreationDate` FROM Lists INNER JOIN Subscriptions ON Lists.ID = Subscriptions.ListID WHERE `UserID`=? AND `GroupID`=?", usrID, grpID)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetUserLists_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	lists := make([]List, 0)
	for rows.Next() {
		var lst List
		if err = rows.Scan(&lst.ID, &lst.Name, &lst.Properties, &lst.Parent, &lst.LatestInvocation, &lst.CreationDate); err != nil {
			db.AddLogEvent(Log{Event: "GetUserLists_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			lists = append(lists, lst)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetUserLists_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetUserLists_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return lists, err
}

//GetUserLists returns a database.Subscription slice containing all the list a user is subscribed to IN ALL THE GROUPS
func (db *SQLiteDB) GetUserLists(usrID int64) ([]List, error) {

	rows, err := db.Query("SELECT `ID`, `Name`, `Properties`, `Parent`, `LatestInvocation`, `CreationDate` FROM Lists INNER JOIN Subscriptions ON Lists.ID = Subscriptions.ListID WHERE `UserID`=?", usrID)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetUserLists_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	lists := make([]List, 0)
	for rows.Next() {
		var lst List
		if err = rows.Scan(&lst.ID, &lst.Name, &lst.Properties, &lst.Parent, &lst.LatestInvocation, &lst.CreationDate); err != nil {
			db.AddLogEvent(Log{Event: "GetUserLists_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			lists = append(lists, lst)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetUserLists_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetUserLists_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return lists, err
}

//GetList returns a list given its ID
func (db *SQLiteDB) GetList(lstID int64) (Subscription, error) {
	var sub Subscription
	err := db.QueryRow("SELECT `ID`,`ListID`,`UserID` FROM Subscriptions WHERE `ID`=?", lstID).Scan(&sub.ID, &sub.ListID, &sub.UserID)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetList_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return sub, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetList_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return sub, err
	default:
		return sub, err
	}
}
