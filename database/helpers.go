package database

import (
	"fmt"
)

//ExecuteRawSQLQuery executes a RAW statement on the databse.
//VERY DANGEROUS
func (db *SQLiteDB) ExecuteRawSQLQuery(query string) (string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err != nil {
			err = rows.Scan(&id, &name)
			return "", err
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}
	//TODO: Fix
	return "", nil
}
