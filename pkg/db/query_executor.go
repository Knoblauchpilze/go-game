package db

import "github.com/KnoblauchPilze/go-game/pkg/errors"

type QueryExecutor interface {
	RunQuery(qb QueryBuilder) error
	RunQueryAndScanSingleResult(qb QueryBuilder, parser RowParser) error
	RunQueryAndScanAllResults(qb QueryBuilder, parser RowParser) error
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
	rows, err := qe.runQueryAndReturnRows(qb)
	if err != nil {
		return err
	}

	rows.Close()
	return nil
}

func (qe *queryExecutorImpl) RunQueryAndScanSingleResult(qb QueryBuilder, parser RowParser) error {
	rows, err := qe.runQueryAndReturnRows(qb)
	if err != nil {
		return err
	}

	defer rows.Close()

	if err := rows.GetSingleValue(parser); err != nil {
		return errors.WrapCode(err, errors.ErrDbCorruptedData)
	}

	return nil
}

func (qe *queryExecutorImpl) RunQueryAndScanAllResults(qb QueryBuilder, parser RowParser) error {
	rows, err := qe.runQueryAndReturnRows(qb)
	if err != nil {
		return err
	}

	defer rows.Close()

	if err := rows.GetAll(parser); err != nil {
		return errors.WrapCode(err, errors.ErrDbCorruptedData)
	}

	return nil
}

func (qe *queryExecutorImpl) runQueryAndReturnRows(qb QueryBuilder) (Rows, error) {
	query, err := qb.Build()
	if err != nil {
		return nil, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	rows := qe.db.Query(query)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}
