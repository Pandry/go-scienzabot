package database

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The logHelpers.go file focuses on the Log table in the database.
// The log table is supposed to be used to log any kind of error that happens within the bot,
//	especially query errors

//AddLogEvent adds a error event into the database
func (db *SQLiteDB) AddLogEvent(evnt Log) error {
	_, err := db.Exec("INSERT INTO Log (`Event`, `RelatedUser` , `RelatedGroup`,"+
		"`Message`, `UpdateValue` , `Error`, `Severity`) VALUES (?,?,?,?,?,?,?)",
		evnt.Event, evnt.RelatedUserID, evnt.RelatedGroupID,
		evnt.Message, evnt.UpdateValue, evnt.Error, evnt.Severity)
	return err
}
