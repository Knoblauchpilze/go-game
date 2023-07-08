package db

// See here:
type Database interface {
	Connect() error
	Disconnect() error

	Query(query Query) Rows
	Execute(query Query) Result
}
