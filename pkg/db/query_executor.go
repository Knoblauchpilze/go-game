package db

import (
	"context"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type QueryExecutor interface {
	RunQueryAndScanSingleResult(ctx context.Context, qb QueryBuilder, parser RowParser) error
	RunQueryAndScanAllResults(ctx context.Context, qb QueryBuilder, parser RowParser) error
	ExecuteQueryAffectingSingleRow(ctx context.Context, qb QueryBuilder) error
}

type queryExecutorImpl struct {
	db Database
}

func NewQueryExecutor(db Database) QueryExecutor {
	return &queryExecutorImpl{
		db: db,
	}
}

func (qe *queryExecutorImpl) RunQueryAndScanSingleResult(ctx context.Context, qb QueryBuilder, parser RowParser) error {
	rows, err := qe.runQueryAndReturnRows(ctx, qb)
	if err != nil {
		return err
	}

	defer rows.Close()

	if err := rows.GetSingleValue(parser); err != nil {
		return errors.WrapCode(err, errors.ErrDbCorruptedData)
	}

	return nil
}

func (qe *queryExecutorImpl) RunQueryAndScanAllResults(ctx context.Context, qb QueryBuilder, parser RowParser) error {
	rows, err := qe.runQueryAndReturnRows(ctx, qb)
	if err != nil {
		return err
	}

	defer rows.Close()

	if err := rows.GetAll(parser); err != nil {
		return errors.WrapCode(err, errors.ErrDbCorruptedData)
	}

	return nil
}

func (qe *queryExecutorImpl) ExecuteQueryAffectingSingleRow(ctx context.Context, qb QueryBuilder) error {
	res, err := qe.executeQueryAndReturn(ctx, qb)
	if err != nil {
		return err
	}

	if res.AffectedRows() != 1 {
		return errors.NewCode(errors.ErrSqlQueryDidNotAffectSingleRow)
	}

	return nil
}

func (qe *queryExecutorImpl) runQueryAndReturnRows(ctx context.Context, qb QueryBuilder) (Rows, error) {
	query, err := qb.Build()
	if err != nil {
		return nil, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	rows := qe.db.Query(ctx, query)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}

func (qe *queryExecutorImpl) executeQueryAndReturn(ctx context.Context, qb QueryBuilder) (Result, error) {
	query, err := qb.Build()
	if err != nil {
		return nil, errors.WrapCode(err, errors.ErrDbRequestCreationFailed)
	}

	res := qe.db.Execute(ctx, query)
	if err := res.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
