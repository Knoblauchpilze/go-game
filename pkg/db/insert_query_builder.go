package db

import (
	"encoding/json"
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

type Convertible interface {
	Convert() interface{}
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
	b.hasReturnValue = !hasReturnValue
}

func (b *scriptQueryBuilder) SetScript(script string) error {
	if len(script) == 0 {
		return errors.NewCode(errors.ErrInvalidSqlTable)
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
		argStr, err := b.argToStr(arg)
		if err != nil {
			return "", err
		}

		args = append(args, fmt.Sprintf("'%s'", argStr))
	}

	return strings.Join(args, ", "), nil
}

func (b *scriptQueryBuilder) argToStr(arg interface{}) (string, error) {
	var raw []byte
	var out string
	var err error

	if convertible, ok := arg.(Convertible); ok {
		raw, err = json.Marshal(convertible.Convert())
		if err == nil {
			out = string(raw)
		}
	} else if str, ok := arg.(string); ok {
		out = str
	} else {
		raw, err = json.Marshal(arg)
		if err == nil {
			out = string(raw)
		}
	}

	return out, err
}
