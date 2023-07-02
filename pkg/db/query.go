package db

type Query interface {
	Valid() bool
	Verbose() bool
	ToSql() string
}

type queryImpl struct {
	sqlCode string
	verbose bool
}

func (q queryImpl) Valid() bool {
	return len(q.sqlCode) > 0
}

func (q queryImpl) Verbose() bool {
	return q.verbose
}

func (q queryImpl) ToSql() string {
	return q.sqlCode
}
