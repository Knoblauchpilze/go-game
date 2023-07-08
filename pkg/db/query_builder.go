package db

type QueryBuilder interface {
	Build() (Query, error)
}
