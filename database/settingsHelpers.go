package database

//For a more detailed explaination about this code, see botSettingsHelpers.go file, in this directory

//GetSettingValue searches for the value in the database
func (db *SQLiteDB) GetSettingValue(key string, group int) (string, error) {
	//We're prepaing a query to execute later
	stmt, err := db.Prepare("SELECT Value FROM `Settings` WHERE `Key` = ? AND `Group` = ? ")
	if err != nil {
		return "", err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()
	//The setting value will go on this string
	var val string
	//Then we execute the query passing the key to the scan function
	err = stmt.QueryRow(key, group).Scan(&val)
	if err != nil {
		return "", err
	}
	//Finally, we return the result
	return val, nil
}

//SetSettingValue sets a value in the bot settings table
func (db *SQLiteDB) SetSettingValue(key string, value string, group int) error {
	stmt, err := db.Prepare("INSERT INTO Settings (`Key`, `Value` , `Group`) VALUES (?,?,?) ON CONFLICT(`Key`, `Group`) DO UPDATE SET `Value` = Excluded.Value")
	if err != nil {
		return err
	}
	defer stmt.Close()

	//And we execute it passing the parameters
	stmt.Exec(key, value, group)

	if err != nil {
		return err
	}

	return err
}
