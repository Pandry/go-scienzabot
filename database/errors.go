package database

//NoRowsAffected is an error thrown when the query did not affect any row
type NoRowsAffected struct {
	error
}
