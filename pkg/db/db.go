package db

// See here:
type Database interface {
	Connect() error
	Disconnect() error

	Query(query Query) QueryRows
	Execute(query Query) Result
}
