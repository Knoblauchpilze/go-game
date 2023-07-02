package db

// See here:
type Database interface {
	Query(query Query) QueryRows
	Execute(query Query) Result
}
