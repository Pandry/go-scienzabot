package database

import (
	"database/sql"
	"errors"
	"fmt"
)

/*
CREATE TABLE IF NOT EXISTS 'Groups' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Title'  TEXT NOT NULL,
	'Ref'	TEXT NOT NULL,
	'Locale'	TEXT DEFAULT ` + DefaultLocale + `,
	'Status'	INTEGER NOT NULL DEFULT 0
);
*/

//AddGroup takes a a database.Group struct as parameter and insert it in the database
//Only ID, Title, Ref and, if present, Locale will be considered since other ones are supposed to be setted later
func (db *SQLiteDB) AddGroup(grp Group) error {
	localeField, localeParameter := "", ""
	if grp.Locale != "" {
		localeField = ", `Locale`"
		localeParameter = ",?"
	}
	query := fmt.Sprintf("INSERT INTO Groups (`ID`, `Title`, `Ref`%s)  VALUES (?,?,?%s)", localeField, localeParameter)

	stmt, err := db.Prepare(query)
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "AddGroupQueryFailed", Message: "Impossible to create the AddGroup preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	var res sql.Result
	if grp.Locale != "" {
		//And we execute it passing the parameters
		res, err = stmt.Exec(grp.ID, grp.Title, grp.Ref, grp.Locale)
	} else {
		res, err = stmt.Exec(grp.ID, grp.Title, grp.Ref)
	}

	if err != nil {
		db.AddLogEvent(Log{Event: "AddGroupQueryFailed", Message: "Impossible to execute the AddGroup preparation query", Error: err.Error()})
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "AddGroupQueryFailedUnknown", Message: "Impossible to execute the AddGroup query due to an undetermined error", Error: err.Error()})
		return err
	}
	if rows > 0 {
		return nil
	}
	db.AddLogEvent(Log{Event: "AddGroupQueryFailedNoAdded", Message: "Impossible to add the AddGroup query: no rows were affected by the query", Error: err.Error()})
	return errors.New("No group was created")

}

//GroupExists returs a bool that indicates if the group exists or not
func (db *SQLiteDB) GroupExists(userID int) bool {
	stmt, err := db.Prepare("SELECT 1 FROM `Groups` WHERE `ID`=?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "GroupExists_QueryFailed", Message: "Impossible to create the GroupExists preparation query", Error: err.Error()})
		return false
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()

	//Then we execute the query passing the userID to the scan function
	qry := stmt.QueryRow(userID)

	var res sql.NullInt64
	err = qry.Scan(&res)
	//And check for errors
	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		return false

	case err != nil:
		db.AddLogEvent(Log{Event: "GroupExists_UnknownError", Message: "Requested a group ID not present in the database but the error is unknown", Error: err.Error()})
		return false

	default:
		//Success
		if res.Valid && res.Int64 == 1 {
			return true
		}
		return false
	}
}

//UpdateDefaultGroupLocale updates the locale of a given group
func (db *SQLiteDB) UpdateDefaultGroupLocale(groupID int, locale string) error {
	if locale == "" {
		return errors.New("EmptyNewLocaleString")
	}

	stmt, err := db.Prepare("UPDATE Groups SET `Locale` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UpdateDefaultLocaleGroup_QueryFailed", Message: "Impossible to create the AddUser preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(locale, groupID)

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "UpdateDefaultLocaleGroup_GroupDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "UpdateDefaultLocaleGroup_GroupDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}
}

//UpdateGroupName updates the name of a given group
func (db *SQLiteDB) UpdateGroupName(groupID int, groupNewName string) error {
	if groupNewName == "" {
		return errors.New("EmptyNewGroupName")
	}

	stmt, err := db.Prepare("UPDATE Groups SET `Name` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UpdateGroupName_QueryFailed", Message: "Impossible to create the UpdateGroupName preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(groupNewName, groupID)

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "UpdateGroupName_GroupDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "UpdateGroupName_GroupDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}
}

//UpdateGroupRef updates the ref of a given group
func (db *SQLiteDB) UpdateGroupRef(groupID int, ref string) error {
	if ref == "" {
		return errors.New("EmptyNewRefString")
	}

	stmt, err := db.Prepare("UPDATE Groups SET `Ref` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UpdateGroupRef_QueryFailed", Message: "Impossible to create the UpdateGroupRef preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(ref, groupID)

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "UpdateGroupRef_GroupDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "UpdateGroupRef_GroupDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}
}

//UpdateGroupTitle updates the title of a given group
func (db *SQLiteDB) UpdateGroupTitle(groupID int, title string) error {
	if title == "" {
		return errors.New("EmptyNewRefString")
	}

	stmt, err := db.Prepare("UPDATE Groups SET `Title` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UpdateGroupTitle_QueryFailed", Message: "Impossible to create the UpdateGroupTitle preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(title, groupID)

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "UpdateGroupTitle_GroupDontExists", Message: "Requested a nickname not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "UpdateGroupTitle_GroupDontExistsUnknown", Message: "Requested a nickname not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}
}

//UpdateGroupStatus updates the status of a given group
func (db *SQLiteDB) UpdateGroupStatus(groupID int, status int) error {
	stmt, err := db.Prepare("UPDATE Groups SET `Status` = ? WHERE `ID` = ?")
	if err != nil {
		//Log the error
		db.AddLogEvent(Log{Event: "UpdateGroupStatus_QueryFailed", Message: "Impossible to create the UpdateGroupStatus preparation query", Error: err.Error()})
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	_, err = stmt.Exec(status, groupID)

	switch {
	case err == sql.ErrNoRows:
		//User does not exist
		db.AddLogEvent(Log{Event: "UpdateGroupStatus_GroupDontExists", Message: "Requested a group not present in the database", Error: err.Error()})
		return err

	case err != nil:
		db.AddLogEvent(Log{Event: "UpdateGroupStatus_GroupDontExistsUnknown", Message: "Requested a group not present in the database but the error is unknown", Error: err.Error()})
		return err

	default:
		//Success
		return nil
	}
}
