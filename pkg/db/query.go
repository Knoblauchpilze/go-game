package db

type Query interface {
	Valid() bool
	ToSql() string
	Verbose() bool
}

type queryImpl struct {
	sqlCode string
	verbose bool
}

func (q queryImpl) Valid() bool {
	return len(q.sqlCode) > 0
}

func (q queryImpl) ToSql() string {
	return q.sqlCode
}

func (q queryImpl) Verbose() bool {
	return q.verbose
}
