package db

type Filter interface {
	Valid() bool
	ToSql() string
}

type filterImpl struct {
	sqlCode string
}

func (f filterImpl) Valid() bool {
	return len(f.sqlCode) > 0
}

func (f filterImpl) ToSql() string {
	return f.sqlCode
}
