package db

type Query interface {
	Valid() bool
	Verbose() bool
	ToSql() (string, error)
}
