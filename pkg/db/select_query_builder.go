package db

import (
	"fmt"
	"strings"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type SelectQueryBuilder interface {
	QueryBuilder

	SetTable(table string) error
	AddProp(prop string) error
	SetVerbose(verbose bool)
}

type selectQueryBuilder struct {
	props   map[string]bool
	table   string
	verbose bool
}

func NewSelectQueryBuilder() SelectQueryBuilder {
	return &selectQueryBuilder{
		props: make(map[string]bool),
	}
}

func (b *selectQueryBuilder) SetTable(table string) error {
	if len(table) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlTable)
	}

	b.table = table
	return nil
}

func (b *selectQueryBuilder) AddProp(prop string) error {
	if len(prop) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlProp)
	}

	if _, ok := b.props[prop]; ok {
		return errors.NewCode(errors.ErrDuplicatedSqlProp)
	}

	b.props[prop] = true
	return nil
}

func (b *selectQueryBuilder) SetVerbose(verbose bool) {
	b.verbose = verbose
}

func (b *selectQueryBuilder) Build() (Query, error) {
	if len(b.table) == 0 {
		return queryImpl{}, errors.WrapCode(errors.NewCode(errors.ErrInvalidSqlTable), errors.ErrSqlTranslationFailed)
	}
	if len(b.props) == 0 {
		return queryImpl{}, errors.WrapCode(errors.NewCode(errors.ErrNoPropInSqlSelectQuery), errors.ErrSqlTranslationFailed)
	}

	propsAsStr := b.propsToStr()
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s", propsAsStr, b.table)

	query := queryImpl{
		sqlCode: sqlQuery,
		verbose: b.verbose,
	}

	return query, nil
}

func (b *selectQueryBuilder) propsToStr() string {
	props := make([]string, 0, len(b.props))
	for prop := range b.props {
		props = append(props, prop)
	}

	return strings.Join(props, ", ")
}
