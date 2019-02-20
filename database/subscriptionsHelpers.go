package database

import (
	"database/sql"
	"errors"
)

/*
type Subscription struct {
	ID     int64
	ListID int64
	UserID int64
}
*/

//AddSubscription adds a subscription. takes as a parameter the userID and the listID
func (db *SQLiteDB) AddSubscription(userID int, listID int) error {
	query, err := db.Exec("INSERT INTO Subscripions (`ListID`, `UserID`) VALUES (?,?)",
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
		db.AddLogEvent(Log{Event: "AddSubscription_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//RemoveSubscription deletes a subscription given its ID
func (db *SQLiteDB) RemoveSubscription(subID int) error {
	query, err := db.Exec("DELETE FROM Subscripions WHERE ID = ?", subID)
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
		db.AddLogEvent(Log{Event: "RemoveSubscription_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
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
	if !rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GellAllSubscriptions_RowsNotFetched", Message: "Some rows in the query were not fetched", Error: err.Error()})
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
	if !rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_RowsNotFetched", Message: "Some rows in the query were not fetched", Error: err.Error()})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return subs, err
}

//GetListSubscribers returns a slice containing all the users in a given list
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
	if !rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetListSubscribers_RowsNotFetched", Message: "Some rows in the query were not fetched", Error: err.Error()})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetListSubscribers_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return subs, err
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
