package db

import (
	"fmt"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type DeleteQueryBuilder interface {
	QueryBuilder

	SetTable(table string) error
	SetFilter(filter Filter) error
	SetVerbose(verbose bool)
}

type deleteQueryBuilder struct {
	table   string
	filter  Filter
	verbose bool
}

func NewDeleteQueryBuilder() DeleteQueryBuilder {
	return &deleteQueryBuilder{}
}

func (b *deleteQueryBuilder) SetTable(table string) error {
	if len(table) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlTable)
	}

	b.table = table
	return nil
}

func (b *deleteQueryBuilder) SetFilter(filter Filter) error {
	if !filter.Valid() {
		return errors.NewCode(errors.ErrInvalidSqlFilter)
	}

	b.filter = filter
	return nil
}

func (b *deleteQueryBuilder) SetVerbose(verbose bool) {
	b.verbose = verbose
}

func (b *deleteQueryBuilder) Build() (Query, error) {
	if len(b.table) == 0 {
		return queryImpl{}, errors.WrapCode(errors.NewCode(errors.ErrInvalidSqlTable), errors.ErrSqlTranslationFailed)
	}

	// https://www.w3schools.com/sql/sql_delete.asp
	sqlQuery := fmt.Sprintf("DELETE FROM %s", b.table)
	if b.filter != nil {
		sqlQuery += fmt.Sprintf(" WHERE %s", b.filter.ToSql())
	}

	query := queryImpl{
		sqlCode: sqlQuery,
		verbose: b.verbose,
	}

	return query, nil
}
