package database

import (
	"strconv"
)

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
func (db *SQLiteDB) ExecuteRawSQLQuery(queryString string) string {

	query, err := db.Exec(queryString)
	if err != nil {
		return "err executing"
	}
	rows, err := query.RowsAffected()
	if err != nil {
		return "err getting affected rows"
	}
	if rows < 1 {
		return "no rows affected"
	}
	return "✅ OK"

}

//QueryRawSQLQuery executes a RAW statement on the databse and return a string containing the result
//VERY DANGEROUS
func (db *SQLiteDB) QueryRawSQLQuery(queryString string) string {

	rows, err := db.Query(queryString)
	defer rows.Close()
	//bkms := make([]Bookmark, 0)
	if err != nil {
		return "Error executing the query"
	}
	res := ""

	var result [][]string
	colTypes, _ := rows.ColumnTypes()
	pointers := make([]interface{}, len(colTypes))
	container := make([]string, len(colTypes))

	for i := range pointers {
		pointers[i] = &container[i]
	}
	i := 0
	for rows.Next() {
		i++
		res += "Rows[" + strconv.Itoa(i) + "]: \n"

		rows.Scan(pointers...)
		result = append(result, container)

		{
			for i, c := range container {
				res += colTypes[i].Name() + ": " + c + "\n"
			}
		}
		res += "\n"
	}
	if rows.NextResultSet() {
		res += "Some roes weren't fetched \n"
	} else if err := rows.Err(); err != nil {
		res += "Error: " + err.Error() + " \n"
	}

	return res

}
