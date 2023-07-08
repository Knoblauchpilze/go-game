package db

import "github.com/KnoblauchPilze/go-game/pkg/errors"

type QueryExecutor interface {
	RunQuery(qb QueryBuilder) error
}

type queryExecutorImpl struct {
	db Database
}

func NewQueryExecutor(db Database) QueryExecutor {
	return &queryExecutorImpl{
		db: db,
	}
}

func (qe *queryExecutorImpl) RunQuery(qb QueryBuilder) error {
	query, err := qb.Build()
	if err != nil {
		return errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	rows := qe.db.Query(query)
	if err := rows.Err(); err != nil {
		return err
	}
	rows.Close()

	return nil
}
