package database

//AddLogEvent adds a error event into the database
func (db *SQLiteDB) AddLogEvent(evnt Log) error {
	_, err := db.Exec("INSERT INTO Log (`Event`, `RelatedUser` , `RelatedGroup`,"+
		"`Message`, `UpdateValue` , `Error`, `Severity`) VALUES (?,?,?,?,?,?,?)",
		evnt.Event, evnt.RelatedUserID, evnt.RelatedGroupID,
		evnt.Message, evnt.UpdateValue, evnt.Error, evnt.Severity)
	return err
}
