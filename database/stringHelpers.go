package database

/*

import (
	"database/sql"
	"log"
)

//TODO: return here when basic queries will be done

//For a more detailed explaination about this code, see botSettingsHelpers.go file, in this directory

//String is referred as a string like the welcomemessage, so custom from a group to another

//GetStringValue searches for the string value in the database
func (db *SQLiteDB) GetStringValue(key string, group int, locale string) (string, error) {
	//We're prepaing a query to execute later
	stmt, err := db.Prepare("SELECT Value FROM `Strings` WHERE `Key` = ? AND `Group` = ? AND `Locale` = ? ")
	if err != nil {
		return "", err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()
	//The setting value will go on this string
	var val string
	//Then we execute the query passing the key to the scan function
	err = stmt.QueryRow(key, group, locale).Scan(&val)
	if err != nil {
		return "", err
	}
	//Finally, we return the result
	return val, nil
}

//GetDefaultStringValue searches for the string value in the database using first the default
//locale of the group, then the  bot default locale, if nothing found tries without specifying the locale
//and if there's still no result, error is returned
func (db *SQLiteDB) GetDefaultStringValue(key string, group int) (string, error) {

	//We're prepaing a query to execute later
	stmt, err := db.Prepare("SELECT Value FROM `Strings` WHERE `Key` = ? AND `Group` = ? AND `Locale` = ? ")
	if err != nil {
		return "", err
	}
	//We want to close the connection to the database once we stop using it
	defer stmt.Close()
	//The setting value will go on this string
	var val string
	//Then we execute the query passing the key to the scan function
	err = stmt.QueryRow(key, group, locale).Scan(&val)
	switch {

	case err == sql.ErrNoRows:
		log.Printf("No user with id %d", id)

	case err != nil:
		log.Fatal(err)

	default:
		//Success

	}
	if err != nil {
		return "", err
	}
	//Finally, we return the result
	return val, nil
}

//SetStringValue sets a value in the bot settings table
func (db *SQLiteDB) SetStringValue(key string, value string, group int, locale string) error {
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
*/
