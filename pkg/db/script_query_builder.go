package db

import (
	"fmt"
	"strings"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type ScriptQueryBuilder interface {
	QueryBuilder

	SetHasReturnValue(hasReturnValue bool)
	SetScript(script string) error
	AddArg(arg interface{}) error
	SetVerbose(verbose bool)
}

type scriptQueryBuilder struct {
	hasReturnValue bool
	script         string
	args           []interface{}
	verbose        bool
}

func NewScriptQueryBuilder() ScriptQueryBuilder {
	return &scriptQueryBuilder{}
}

func (b *scriptQueryBuilder) SetHasReturnValue(hasReturnValue bool) {
	b.hasReturnValue = hasReturnValue
}

func (b *scriptQueryBuilder) SetScript(script string) error {
	if len(script) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlScript)
	}

	b.script = script
	return nil
}

func (b *scriptQueryBuilder) AddArg(arg interface{}) error {
	if arg == nil {
		return errors.NewCode(errors.ErrInvalidSqlScriptArg)
	}

	b.args = append(b.args, arg)

	return nil
}

func (b *scriptQueryBuilder) SetVerbose(verbose bool) {
	b.verbose = verbose
}

func (b *scriptQueryBuilder) Build() (Query, error) {
	if len(b.script) == 0 {
		return queryImpl{}, errors.WrapCode(errors.NewCode(errors.ErrInvalidSqlScript), errors.ErrSqlTranslationFailed)
	}

	argsAsStr, err := b.argsToStr()
	if err != nil {
		return queryImpl{}, errors.WrapCode(err, errors.ErrSqlTranslationFailed)
	}

	query := queryImpl{
		verbose: b.verbose,
	}

	if b.hasReturnValue {
		query.sqlCode = fmt.Sprintf("SELECT * FROM %s(%s)", b.script, argsAsStr)
	} else {
		query.sqlCode = fmt.Sprintf("SELECT %s(%s)", b.script, argsAsStr)
	}

	return query, nil
}

func (b *scriptQueryBuilder) argsToStr() (string, error) {
	args := make([]string, 0, len(b.args))
	for _, arg := range b.args {
		argStr, err := argToStr(arg)
		if err != nil {
			return "", err
		}

		args = append(args, fmt.Sprintf("'%s'", argStr))
	}

	return strings.Join(args, ", "), nil
}
