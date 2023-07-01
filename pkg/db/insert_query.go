package db

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Convertible interface {
	Convert() interface{}
}

type insertQuery struct {
	script     string
	args       []interface{}
	skipReturn bool
	verbose    bool
}

func (q insertQuery) Valid() bool {
	return len(q.script) > 0
}

func (q insertQuery) Verbose() bool {
	return q.verbose
}

func (q insertQuery) ToSql() (string, error) {
	argsAsStr := make([]string, 0)

	for _, arg := range q.args {
		cvrt, ok := arg.(Convertible)

		var raw []byte
		var err error

		if ok {
			raw, err = json.Marshal(cvrt.Convert())
		} else {
			str, ok := arg.(string)

			if ok {
				raw = []byte(str)
			} else {
				raw, err = json.Marshal(arg)
			}
		}

		if err != nil {
			return "", err
		}

		argAsStr := fmt.Sprintf("'%s'", string(raw))

		argsAsStr = append(argsAsStr, argAsStr)
	}

	var query string

	switch q.skipReturn {
	case false:
		query = fmt.Sprintf("SELECT * from %s(%s)", q.script, strings.Join(argsAsStr, ", "))
	default:
		query = fmt.Sprintf("SELECT %s(%s)", q.script, strings.Join(argsAsStr, ", "))
	}

	return query, nil
}
