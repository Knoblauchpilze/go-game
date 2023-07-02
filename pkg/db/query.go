package db

import "github.com/KnoblauchPilze/go-game/pkg/logger"

type Query interface {
	Valid() bool
	ToSql() string
}

type queryImpl struct {
	sqlCode string
	verbose bool
}

func (q queryImpl) Valid() bool {
	return len(q.sqlCode) > 0
}

func (q queryImpl) ToSql() string {
	if q.verbose {
		logger.Tracef("sql query: %s", q.sqlCode)
	}
	return q.sqlCode
}
