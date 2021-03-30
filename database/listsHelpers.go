package database

import (
	"database/sql"
	"errors"
	"scienzabot/consts"
	"scienzabot/utils"
	"time"
)

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The listsHelpers.go file focuses on the Lists table in the database.
// The lists is supposed to keep track of the lists the bot is in charge to manage and offers
//	functions to do so

//TODO: Refactor methods

//GetAvailableLists returns an array of lists a user can subscribe to
func (db *SQLiteDB) GetAvailableLists(groupID int64, userID int, limit int, offset int) ([]List, error) {
	var rows *sql.Rows
	var err error
	if limit < 1 || offset < 0 {
		rows, err = db.Query("SELECT `ID`, `Name`, `Properties` FROM Lists WHERE `GroupID` = ? AND `ID` NOT IN (SELECT `ListID` FROM Subscriptions WHERE `UserID` = ?)", groupID, userID)
	} else {
		rows, err = db.Query("SELECT `ID`, `Name`, `Properties` FROM Lists WHERE `GroupID` = ? AND `ID` NOT IN (SELECT `ListID` FROM Subscriptions WHERE `UserID` = ?) LIMIT ? OFFSET ?", groupID, userID, limit, offset)
	}

	defer rows.Close()
	resultLists := make([]List, 0)

	if err != nil {
		db.AddLogEvent(Log{Event: "GetAvailableLists_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return resultLists, err
	}

	for rows.Next() {
		var (
			ID, prop int64
			Name     string
		)

		if err = rows.Scan(&ID, &Name, &prop); err != nil {
			db.AddLogEvent(Log{Event: "GetAvailableLists_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			resultLists = append(resultLists, List{ID: ID, Name: Name, GroupID: groupID, Properties: prop})
		}
	}

	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetAvailableLists_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetAvailableLists_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return resultLists, err
}

//GetLists returns an array of lists given a group
func (db *SQLiteDB) GetLists(groupID int64) ([]List, error) {
	rows, err := db.Query("SELECT `ID`, `Name`, `Properties`, `LatestInvocation` FROM Lists WHERE `GroupID` = ?", groupID)
	defer rows.Close()
	resultLists := make([]List, 0)
	if err != nil {
		db.AddLogEvent(Log{Event: "GetLists_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return resultLists, err
	}
	for rows.Next() {
		var (
			inv     sql.NullString
			tmpList List
		)

		if err = rows.Scan(&tmpList.ID, &tmpList.Name, &tmpList.Properties, &inv); err != nil {
			db.AddLogEvent(Log{Event: "GetLists_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			tmpList.LatestInvocation, _ = time.Parse(consts.TimeFormatString, inv.String)
			resultLists = append(resultLists, tmpList)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetLists_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetLists_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return resultLists, err
}

//GetList returns a list given its ID
func (db *SQLiteDB) GetList(listID int64) (List, error) {
	l := List{}
	var li sql.NullString
	var p sql.NullInt64
	err := db.QueryRow("SELECT `Name`,`GroupID`, `Properties`, `CreationDate`, `LatestInvocation`, `Parent` FROM Lists WHERE `ID` = ? LIMIT 1", listID).
		Scan(&l.Name, &l.GroupID, &l.Properties, &l.CreationDate, &li, &p)
		//TODO: Handle the error
	l.LatestInvocation, _ = time.Parse(consts.TimeFormatString, li.String)

	l.Parent = p.Int64
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetList_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return l, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetList_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return l, err
	default:
		return l, nil
	}

	/*resultLists := make([]List, 0)
	if err != nil {
		db.AddLogEvent(Log{Event: "GetLists_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return resultLists, err
	}
	for rows.Next() {
		var (
			inv     sql.NullString
			tmpList List
		)

		if err = rows.Scan(&tmpList.ID, &tmpList.Name, &tmpList.Properties, &inv); err != nil {
			db.AddLogEvent(Log{Event: "GetLists_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			tmpList.LatestInvocation, _ = time.Parse(consts.TimeFormatString, inv.String)
			resultLists = append(resultLists, tmpList)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetLists_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetLists_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return resultLists, err*/
}

//AddList takes a a database.List struct as parameter and insert it in the database
func (db *SQLiteDB) AddList(lst List) error {
	//lst.Name, lst.GroupID, lst.GroupIndipendent, lst.InviteOnly
	stmt, err := db.Prepare("INSERT INTO Lists (`Name`, `GroupID`, `Properties`)  VALUES (?,?,?)")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "AddList_QueryFailed", Message: "Impossible to create the AddList preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(lst.Name, lst.GroupID, lst.Properties)

	if err != nil {
		db.AddLogEvent(Log{Event: "AddList_QueryFailed", Message: "Impossible to execute the AddList preparation query", Error: err.Error()})
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "AddList_QueryFailedUnknown", Message: "Impossible to execute the AddList query due to an undetermined error", Error: err.Error()})
		return err
	}
	if rows > 0 {
		return nil
	}
	db.AddLogEvent(Log{Event: "AddList_QueryFailedNoAdded", Message: "Impossible to execute the AddList query: no rows were affected by the query", Error: err.Error()})
	return errors.New("No list was created")
}

//DeleteList deletes a list from the database given its listID
func (db *SQLiteDB) DeleteList(listID int) error {
	stmt, err := db.Prepare("DELETE FROM Lists WHERE ID = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "DeleteList_QueryFailed", Message: "Impossible to create the DeleteList preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(listID)

	if err != nil {
		db.AddLogEvent(Log{Event: "DeleteList_QueryFailed", Message: "Impossible to execute the DeleteList preparation query", Error: err.Error()})
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "DeleteList_QueryFailedUnknown", Message: "Impossible to execute the DeleteList query due to an undetermined error", Error: err.Error()})
		return err
	}
	if rows > 0 {
		return nil
	}
	db.AddLogEvent(Log{Event: "DeleteList_QueryFailedNoAdded", Message: "Impossible to execute the DeleteList query: no rows were affected by the query", Error: err.Error()})
	return NoRowsAffected{error: errors.New("No list was deleted")}
}

//DeleteListByName deletes a list from the database given its name and the group it belongs to
func (db *SQLiteDB) DeleteListByName(groupID int64, listName string) error {
	query, err := db.Exec("DELETE FROM Lists WHERE `GroupID` = ? AND `Name` = ?", groupID, listName)
	if err != nil {
		db.AddLogEvent(Log{Event: "DeleteListByName_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "DeleteListByName_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "DeleteListByName_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

/*
CREATE TABLE IF NOT EXISTS 'Lists' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'Name'  TEXT NOT NULL UNIQUE,
	'Group'	INTEGER NOT NULL,
	'GroupIndipendent'  INTEGER DEFAULT 0,
	'InviteOnly'  INTEGER DEFAULT 0,
	FOREIGN KEY('Group') REFERENCES Groups('ID')
);
*/

//RenameList takes a a database.List struct as parameter and insert it in the database
func (db *SQLiteDB) RenameList(listID int, newListName string) error {
	//lst.Name, lst.GroupID, lst.GroupIndipendent, lst.InviteOnly
	stmt, err := db.Prepare("UPDATE Lists SET `Name`=? WHERE ID = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "RenameList_QueryFailed", Message: "Impossible to create the RenameList preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(newListName, listID)

	if err != nil {
		db.AddLogEvent(Log{Event: "RenameList_QueryFailed", Message: "Impossible to execute the RenameList preparation query", Error: err.Error()})
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "RenameList_QueryFailedUnknown", Message: "Impossible to execute the RenameList query due to an undetermined error", Error: err.Error()})
		return err
	}
	if rows > 0 {
		return nil
	}
	db.AddLogEvent(Log{Event: "RenameList_QueryFailedNoAdded", Message: "Impossible to execute the RenameList query: no rows were affected by the query", Error: err.Error()})
	return errors.New("No list was renamed")
}

//EditList takes a a database.List struct as parameter and insert it in the database
func (db *SQLiteDB) EditList(listID int, newListValues List) error {
	//lst.Name, lst.GroupID, lst.GroupIndipendent, lst.InviteOnly
	stmt, err := db.Prepare("UPDATE Lists SET `Name`=?, `GroupID`=?, `Properties`=?, `Status`=? WHERE ID = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "EditList_QueryFailed", Message: "Impossible to create the EditList preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(newListValues.Name, newListValues.GroupID, newListValues.Properties, listID)

	if err != nil {
		db.AddLogEvent(Log{Event: "EditList_QueryFailed", Message: "Impossible to execute the EditList preparation query", Error: err.Error()})
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "EditList_QueryFailedUnknown", Message: "Impossible to execute the EditList query due to an undetermined error", Error: err.Error()})
		return err
	}
	if rows > 0 {
		return nil
	}
	db.AddLogEvent(Log{Event: "EditList_QueryFailedNoAdded", Message: "Impossible to execute the EditList query: no rows were affected by the query", Error: err.Error()})
	return errors.New("No list was edited")
}

//propertyReturnValue is basically a int, but is only returned from GetListProperties and is used to "concatenate" the output of the
//  function to some other method, like db.GetListProperties(groupID).ListIsGroupIndipendent()
type propertyReturnValue struct {
	int
}

//GetListProperties gets the properties given a groupID
func (db *SQLiteDB) GetListProperties(groupID int64) propertyReturnValue {

	stmt, err := db.Prepare("SELECT `Properties` FROM Lists WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetListProperties_QueryFailed", Message: "Impossible to create the GetListProperties preparation query", Error: err.Error()})
		return propertyReturnValue{-1}
	}
	defer stmt.Close()

	res, err := stmt.Query(groupID)
	if err != nil {
		db.AddLogEvent(Log{Event: "GetListProperties_QueryExecutionFailed", Message: "Impossible to execute the GetListProperties preparation query", RelatedGroupID: groupID, Error: err.Error()})
		return propertyReturnValue{0}
	}
	defer res.Close()

	var prop sql.NullInt64

	for res.Next() {
		err = res.Scan(&prop)

		if err != nil {
			db.AddLogEvent(Log{Event: "GetListProperties_GroupDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", RelatedGroupID: groupID, Error: err.Error()})
			continue
		}

	}

	return propertyReturnValue{int(prop.Int64)}
}

//UpdateListLastInvokation takes in iput a list id and its latest invokation and update the invokation in the database
func (db *SQLiteDB) UpdateListLastInvokation(listID int64, invTime time.Time) error {
	query, err := db.Exec("UPDATE Lists SET `LatestInvocation`=? WHERE ID = ?", invTime.Format(consts.TimeFormatString), listID)
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateListLastInvokation_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateListLastInvokation_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "UpdateListLastInvokation_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//ListIsGroupIndipendent returns a bool that indicates if the list is valid in all the groups
func (p propertyReturnValue) ListIsGroupIndipendent() bool {
	return utils.HasPermission(p.int, consts.ListPropertyGroupIndipendent)
}

//ListIsInviteOnly returns a bool that indicates if the list is valid in all the groups
func (p propertyReturnValue) ListIsInviteOnly() bool {
	return utils.HasPermission(p.int, consts.ListPropertyGroupLocked)
}
