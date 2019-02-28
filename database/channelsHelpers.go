package database

import (
	"database/sql"
	"errors"
)

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The channelsHelpers.go file focuses on the Channels table in the database.
// The channel table is supposed to support a feature is basically a feature that permits to forward a discussion to a channel
// It isn't implemented yet.

//AddChannel inserts a new channel in the database
func (db *SQLiteDB) AddChannel(chn Channel) error {
	query, err := db.Exec("INSERT INTO Channels (`ID`, `GroupID`, `Name`, `Ref`) VALUES (?,?,?,?)",
		chn.ID, chn.GroupID, chn.Name, chn.Ref)
	if err != nil {
		db.AddLogEvent(Log{Event: "CreateChannel_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "CreateChannel_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "CreateChannel_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//RemoveChannel removes a channel form the database given its ID
func (db *SQLiteDB) RemoveChannel(chID int64) error {
	query, err := db.Exec("DELETE FROM Channels WHERE `ID`=?",
		chID)
	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveChannel_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "RemoveChannel_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "RemoveChannel_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//UpdateChannel update a database from the database
//It replaces all the values
//The fist paramter is the ID of the channel
//The second parameter is a new channel struct
func (db *SQLiteDB) UpdateChannel(chnID int64, chn Channel) error {
	if chn.ID == 0 {
		chn.ID = chnID
	}
	query, err := db.Exec("UPDATE Channels SET `ID`=?, `GroupID`=?, `Name`=?, `Ref`=? WHERE `ID`=?",
		chn.ID, chn.GroupID, chn.Name, chn.Ref, chnID)
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateChannel_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "UpdateChannel_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "UpdateChannel_NoRowsAffected", Message: "No rows affected"})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//GetChannel returns a single channel given its ID
func (db *SQLiteDB) GetChannel(chnID int64) (Channel, error) {
	var chn Channel
	err := db.QueryRow("SELECT `ID`, `GroupID`, `Name`, `Ref` FROM Channels WHERE `ID`=?",
		chnID).Scan(&chn.ID, &chn.GroupID, &chn.Name, &chn.Ref)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetChannel_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return chn, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetChannel_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return chn, err
	default:
		return chn, nil
	}
}

//GetAllChannels returns all the channels in the database
func (db *SQLiteDB) GetAllChannels() ([]Channel, error) {
	rows, err := db.Query("SELECT `ID`, `GroupID`, `Name`, `Ref` FROM Channels")
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetAllChannels_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	chns := make([]Channel, 0)
	for rows.Next() {
		var chn Channel
		if err = rows.Scan(&chn.ID, &chn.GroupID, &chn.Name, &chn.Ref); err != nil {
			db.AddLogEvent(Log{Event: "GetAllChannels_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			chns = append(chns, chn)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetAllChannels_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	}
	if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetAllChannels_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}

	return chns, err
}

//GetChannelsByName returns all the channels in the database whose name is given in input
func (db *SQLiteDB) GetChannelsByName(qry string) ([]Channel, error) {
	rows, err := db.Query("SELECT `ID`, `GroupID`, `Name`, `Ref` FROM Channels WHERE Name LIKE '%?%'", qry)
	defer rows.Close()
	if err != nil {
		db.AddLogEvent(Log{Event: "GetChannelsByName_ErorExecutingTheQuery", Message: "Impossible to get afftected rows", Error: err.Error()})
		return nil, err
	}
	chns := make([]Channel, 0)
	for rows.Next() {
		var chn Channel
		if err = rows.Scan(&chn.ID, &chn.GroupID, &chn.Name, &chn.Ref); err != nil {
			db.AddLogEvent(Log{Event: "GetChannelsByName_RowQueryFetchResultFailed", Message: "Impossible to get data from the row", Error: err.Error()})
		} else {
			chns = append(chns, chn)
		}
	}
	if rows.NextResultSet() {
		db.AddLogEvent(Log{Event: "GetChannelsByName_RowsNotFetched", Message: "Some rows in the query were not fetched"})
	} else if err := rows.Err(); err != nil {
		db.AddLogEvent(Log{Event: "GetChannelsByName_UnknowQueryError", Message: "An unknown error was thrown", Error: err.Error()})
	}
	return chns, err
}
