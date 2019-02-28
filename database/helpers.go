package database

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The elpers.go file is a generic file where it's possible to find functions that does not belong to
//	any table or feature in particular

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
		//fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}
	//TODO: Fix
	return "", nil
}
