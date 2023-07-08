package db

import (
	"fmt"
	"strings"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type InsertQueryBuilder interface {
	QueryBuilder

	SetTable(table string) error
	AddElement(column string, value interface{}) error
	SetVerbose(verbose bool)
}

type insertQueryBuilder struct {
	columns map[string]bool
	props   []sqlProp
	table   string
	verbose bool
}

type sqlProp struct {
	column string
	value  interface{}
}

func NewInsertQueryBuilder() InsertQueryBuilder {
	return &insertQueryBuilder{
		columns: make(map[string]bool),
	}
}

func (b *insertQueryBuilder) SetTable(table string) error {
	if len(table) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlTable)
	}

	b.table = table
	return nil
}

func (b *insertQueryBuilder) AddElement(column string, value interface{}) error {
	if len(column) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlColumn)
	}

	if _, ok := b.columns[column]; ok {
		return errors.NewCode(errors.ErrDuplicatedSqlColumn)
	}

	prop := sqlProp{
		column: column,
		value:  value,
	}
	b.columns[column] = true
	b.props = append(b.props, prop)
	return nil
}

func (b *insertQueryBuilder) SetVerbose(verbose bool) {
	b.verbose = verbose
}

func (b *insertQueryBuilder) Build() (Query, error) {
	if len(b.table) == 0 {
		return queryImpl{}, errors.WrapCode(errors.NewCode(errors.ErrInvalidSqlTable), errors.ErrSqlTranslationFailed)
	}
	if len(b.props) == 0 {
		return queryImpl{}, errors.WrapCode(errors.NewCode(errors.ErrNoColumnInSqlInsertQuery), errors.ErrSqlTranslationFailed)
	}

	columnsAsStr := b.columnsToStr()
	valuesAsStr, err := b.valuesToStr()
	if err != nil {
		return queryImpl{}, errors.WrapCode(err, errors.ErrSqlTranslationFailed)
	}

	// https://www.w3schools.com/sql/sql_insert.asp
	sqlQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", b.table, columnsAsStr, valuesAsStr)

	query := queryImpl{
		sqlCode: sqlQuery,
		verbose: b.verbose,
	}

	return query, nil
}

func (b *insertQueryBuilder) columnsToStr() string {
	var columns []string

	for _, prop := range b.props {
		columns = append(columns, prop.column)
	}

	return strings.Join(columns, ", ")
}

func (b *insertQueryBuilder) valuesToStr() (string, error) {
	var values []string

	for _, prop := range b.props {
		arg, err := argToStr(prop.value)
		if err != nil {
			return "", err
		}

		values = append(values, fmt.Sprintf("'%s'", arg))
	}

	return strings.Join(values, ", "), nil
}
