package database

//AddLogEvent adds a error event into the database
func (db *SQLiteDB) AddLogEvent(evnt Log) error {
	stmt, err := db.Prepare("INSERT INTO Log (`Event`, `RelatedUser` , `RelatedGroup`," +
		"`Message`, `UpdateValue` , `Error`, `Severity`) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	stmt.Exec(evnt.Event, evnt.RelatedUserID, evnt.RelatedGroupID,
		evnt.Message, evnt.UpdateValue, evnt.Error, evnt.Severity)

	return nil
}
