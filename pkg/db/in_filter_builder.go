package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type InFilterBuilder interface {
	FilterBuilder

	SetKey(key string) error
	AddValue(value interface{}) error
}

type inFilterBuilder struct {
	key    string
	values []interface{}
}

func NewInFilterBuilder() InFilterBuilder {
	return &inFilterBuilder{}
}

func (b *inFilterBuilder) SetKey(key string) error {
	if len(key) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlComparisonKey)
	}

	b.key = key
	return nil
}

func (b *inFilterBuilder) AddValue(value interface{}) error {
	if value == nil {
		return errors.NewCode(errors.ErrInvalidSqlComparisonValue)
	}

	b.values = append(b.values, value)

	return nil
}

func (b *inFilterBuilder) Build() (Filter, error) {
	if len(b.key) == 0 {
		return filterImpl{}, errors.WrapCode(errors.NewCode(errors.ErrInvalidSqlComparisonKey), errors.ErrSqlTranslationFailed)
	}
	if len(b.values) == 0 {
		return filterImpl{}, errors.WrapCode(errors.NewCode(errors.ErrNoValuesInSqlComparison), errors.ErrSqlTranslationFailed)
	}

	valuesAsStr := b.valuesToStr()
	sqlFilter := fmt.Sprintf("%s in (%s)", b.key, valuesAsStr)

	filter := filterImpl{
		sqlCode: sqlFilter,
	}

	return filter, nil
}

func (b *inFilterBuilder) valuesToStr() string {
	values := make([]string, 0, len(b.values))
	for _, value := range b.values {
		valueStr := valueToStr(value)
		values = append(values, fmt.Sprintf("'%s'", valueStr))
	}

	return strings.Join(values, ", ")
}

func valueToStr(value interface{}) string {
	if t, ok := value.(time.Time); ok {
		return fmt.Sprintf("%v", t.Format(time.RFC3339))
	}

	return fmt.Sprintf("%v", value)
}
