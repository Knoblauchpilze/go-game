package db

type FilterBuilder interface {
	Build() (Filter, error)
}
