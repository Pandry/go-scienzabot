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
	stmt, err := db.Prepare("INSERT INTO Subscripions (`ListID`, `UserID`) VALUES (?,?)")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "AddSubscription_QueryFailed", Message: "Impossible to create the AddSubscription preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec()
	if err != nil {
		db.AddLogEvent(Log{Event: "AddSubscription_ExecutionQueryFailed", Message: "Impossible to execute the AddSubscription query", Error: err.Error()})
	}
	return err
}

//RemoveSubscription deletes a subscription given its ID
func (db *SQLiteDB) RemoveSubscription(subID int) error {
	stmt, err := db.Prepare("DELETE FROM Subscripions WHERE ID = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "RemoveSubscription_QueryFailed", Message: "Impossible to create the RemoveSubscription preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(subID)

	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveSubscription_QueryFailed", Message: "Impossible to execute the RemoveSubscription preparation query", Error: err.Error()})
		return err
	}
	rows, err := res.RowsAffected()

	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveSubscription_QueryFailedUnknown", Message: "Impossible to execute the RemoveSubscription query due to an undetermined error", Error: err.Error()})
		return err
	}

	if rows > 0 {
		return nil
	}
	db.AddLogEvent(Log{Event: "RemoveSubscription_QueryFailedNoAdded", Message: "Impossible to execute the RemoveSubscription query: no rows were affected by the query", Error: err.Error()})
	return NoRowsAffected{error: errors.New("No subscription was deleted")}
}

//GellAllSubscriptions returns a slice containing all the users in a given list
func (db *SQLiteDB) GellAllSubscriptions() ([]Subscription, error) {
	var result []Subscription

	stmt, err := db.Prepare("SELECT `ID`,`ListID`,`UserID` FROM Subscriptions")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GellAllSubscriptions_QueryFailed", Message: "Impossible to create the GellAllSubscriptions preparation query", Error: err.Error()})
		return nil, err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	query, err := stmt.Query()

	if err != nil {
		db.AddLogEvent(Log{Event: "GellAllSubscriptions_QueryExecutionFailed", Message: "Impossible execute the GellAllSubscriptions query", Error: err.Error()})
		return nil, err
	}

	var id, listID, userID int64

	for query.Next() {
		var sub Subscription
		//We then scan the query
		err = query.Scan(&id, &listID, &userID)
		sub.ID = id
		sub.ListID = listID
		sub.UserID = userID
		//And check for errors
		switch {
		case err == sql.ErrNoRows:
			//Group does not exist - DAFUQ?!
			db.AddLogEvent(Log{Event: "GellAllSubscriptions_NoRows", Message: "Subscription not found in database - this should NEVER happen, something's wrong!", Error: err.Error()})
			//return result, err

		case err != nil:
			db.AddLogEvent(Log{Event: "GellAllSubscriptions_UncaughtError", Message: "A error happened and it was not identified", Error: err.Error()})
			//return result, err

		default:
			//Success
			result = append(result, sub)
		}
	}
	return result, err
}

//GetSubscribedUsers returns a slice containing all the users in a given list
func (db *SQLiteDB) GetSubscribedUsers(lstID int64) ([]Subscription, error) {
	var result []Subscription

	stmt, err := db.Prepare("SELECT `ID`,`ListID`,`UserID` FROM Subscriptions WHERE `ListID`=?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_QueryFailed", Message: "Impossible to create the GetSubscribedUsers preparation query", Error: err.Error()})
		return nil, err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	query := stmt.QueryRow(lstID)

	if err != nil {
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_QueryExecutionFailed", Message: "Impossible execute the GetSubscribedUsers query", Error: err.Error()})
		return nil, err
	}

	var id, listID, userID int64
	var sub Subscription

	//We then scan the query
	err = query.Scan(&id, &listID, &userID)

	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_NoRows", Message: "Subscription not found in database - this should NEVER happen, something's wrong!", Error: err.Error()})
		return nil, err

	case err != nil:
		db.AddLogEvent(Log{Event: "GetSubscribedUsers_UncaughtError", Message: "A error happened and it was not identified", Error: err.Error()})
		return nil, err

	default:
		sub.ID = id
		sub.ListID = listID
		sub.UserID = userID
		return result, err
	}
}

//GetListSubscribers returns a slice containing all the users in a given list
func (db *SQLiteDB) GetListSubscribers(usrID int64) ([]Subscription, error) {
	var result []Subscription

	stmt, err := db.Prepare("SELECT `ID`,`ListID`,`UserID` FROM Subscriptions WHERE `ID`=?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetListSubscribers_QueryFailed", Message: "Impossible to create the GetListSubscribers preparation query", Error: err.Error()})
		return nil, err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	query := stmt.QueryRow(usrID)

	if err != nil {
		db.AddLogEvent(Log{Event: "GetListSubscribers_QueryExecutionFailed", Message: "Impossible execute the GetListSubscribers query", Error: err.Error()})
		return nil, err
	}

	var id, listID, userID int64
	var sub Subscription

	//We then scan the query
	err = query.Scan(&id, &listID, &userID)

	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetListSubscribers_NoRows", Message: "Subscription not found in database - this should NEVER happen, something's wrong!", Error: err.Error()})
		return nil, err

	case err != nil:
		db.AddLogEvent(Log{Event: "GetListSubscribers_UncaughtError", Message: "A error happened and it was not identified", Error: err.Error()})
		return nil, err

	default:
		sub.ID = id
		sub.ListID = listID
		sub.UserID = userID
		return result, nil
	}
}
