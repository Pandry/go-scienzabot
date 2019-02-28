package database

//The database package is supposed to contain all the database functions and helpers functions
// A helper function is a function that interfaces with the database via a query.
// The helper functions were made to avoid a mantainer to interface directly with the database.
// Each file in the ^([a-zA-Z]+)Helpers.go$ format is supposed to be a "table" helper (Basically
//	a file that have queries about only one table in the database, to keep things tidy.)
// The table name is the $1 group in the above regex.

// The errors.go file contains some custom error the package can "throw" to describe better a situation

//NoRowsAffected is an error thrown when the query did not affect any row
type NoRowsAffected struct {
	error
}

//ParameterError is an error thrown when a parameter misses or is invalid
type ParameterError struct {
	error
}
