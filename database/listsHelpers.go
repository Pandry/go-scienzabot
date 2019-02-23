package database

import (
	"database/sql"
	"errors"
	"scienzabot/consts"
	"scienzabot/utils"
)

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

//TODO: Refactor methods

//GetLists returns an array of lists given a group
func (db *SQLiteDB) GetLists(groupID int64) ([]List, error) {

	var resultLists []List
	//asd.GroupID, asd.GroupIndipendent, asd.ID, asd.InviteOnly, asd.Name

	stmt, err := db.Prepare("SELECT `ID`, `Name`, `Properties` FROM Lists WHERE `GroupID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GetLists_QueryFailed", Message: "Impossible to create the GetLists preparation query",
			RelatedGroupID: groupID, Error: err.Error()})
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Query(groupID)
	if err != nil {
		db.AddLogEvent(Log{Event: "GetLists_QueryExecutionFailed", Message: "Impossible to execute the GetLists preparation query",
			RelatedGroupID: groupID, Error: err.Error()})
		return nil, err
	}
	defer res.Close()
	var ID, prop sql.NullInt64
	var Name sql.NullString
	var tmpList List
	for res.Next() {
		err = res.Scan(&ID, &Name, &prop)

		if err != nil {
			db.AddLogEvent(Log{Event: "GetLists_GroupDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown",
				RelatedGroupID: groupID, Error: err.Error()})
			continue
		}

		tmpList = List{ID: ID.Int64, GroupID: groupID,
			Name: Name.String, Properties: prop.Int64}
		if len(resultLists) == 0 {
			resultLists = []List{tmpList}
		} else {
			resultLists = append(resultLists, tmpList)
		}
	}

	return resultLists, nil
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
	return NoRowsAffected{error: errors.New("No list was created")}
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

//ListIsGroupIndipendent returns a bool that indicates if the list is valid in all the groups
func (p propertyReturnValue) ListIsGroupIndipendent() bool {
	return utils.HasPermission(p.int, consts.ListPropertyGroupIndipendent)
}

//ListIsInviteOnly returns a bool that indicates if the list is valid in all the groups
func (p propertyReturnValue) ListIsInviteOnly() bool {
	return utils.HasPermission(p.int, consts.ListPropertyGroupLocked)
}
