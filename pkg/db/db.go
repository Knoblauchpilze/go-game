package db

// See here:
type Database interface {
	Query(query Query) (Rows, error)
	Execute(query Query) (Result, error)
}
