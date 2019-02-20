package database

//NoRowsAffected is an error thrown when the query did not affect any row
type NoRowsAffected struct {
	error
}

//ParameterError is an error thrown when a parameter misses or is invalid
type ParameterError struct {
	error
}
