package db

import (
	"fmt"
	"strings"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type UpdateQueryBuilder interface {
	QueryBuilder

	SetTable(table string) error
	AddUpdate(column string, newValue interface{}) error
	SetFilter(filter Filter) error
	SetVerbose(verbose bool)
}

type updateQueryBuilder struct {
	columns map[string]bool
	props   []sqlProp
	table   string
	filter  Filter
	verbose bool
}

func NewUpdateQueryBuilder() UpdateQueryBuilder {
	return &updateQueryBuilder{
		columns: make(map[string]bool),
	}
}

func (b *updateQueryBuilder) SetTable(table string) error {
	if len(table) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlTable)
	}

	b.table = table
	return nil
}

func (b *updateQueryBuilder) AddUpdate(column string, newValue interface{}) error {
	if len(column) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlColumn)
	}

	if _, ok := b.columns[column]; ok {
		return errors.NewCode(errors.ErrDuplicatedSqlColumn)
	}

	prop := sqlProp{
		column: column,
		value:  newValue,
	}
	b.columns[column] = true
	b.props = append(b.props, prop)
	return nil
}

func (b *updateQueryBuilder) SetFilter(filter Filter) error {
	if !filter.Valid() {
		return errors.NewCode(errors.ErrInvalidSqlFilter)
	}

	b.filter = filter
	return nil
}

func (b *updateQueryBuilder) SetVerbose(verbose bool) {
	b.verbose = verbose
}

func (b *updateQueryBuilder) Build() (Query, error) {
	if len(b.table) == 0 {
		return queryImpl{}, errors.WrapCode(errors.NewCode(errors.ErrInvalidSqlTable), errors.ErrSqlTranslationFailed)
	}
	if len(b.props) == 0 {
		return queryImpl{}, errors.WrapCode(errors.NewCode(errors.ErrNoColumnInSqlUpdateQuery), errors.ErrSqlTranslationFailed)
	}

	updates, err := b.updatesToStr()
	if err != nil {
		return queryImpl{}, errors.WrapCode(err, errors.ErrSqlTranslationFailed)
	}

	// https://www.w3schools.com/sql/sql_update.asp
	sqlQuery := fmt.Sprintf("UPDATE %s SET %s", b.table, updates)
	if b.filter != nil {
		sqlQuery += fmt.Sprintf(" WHERE %s", b.filter.ToSql())
	}

	query := queryImpl{
		sqlCode: sqlQuery,
		verbose: b.verbose,
	}

	return query, nil
}

func (b *updateQueryBuilder) updatesToStr() (string, error) {
	var updates []string

	for _, prop := range b.props {
		update, err := sqlPropAsUpdateToStr(prop)
		if err != nil {
			return "", err
		}

		updates = append(updates, update)
	}

	return strings.Join(updates, ", "), nil
}
