package db

import (
	"fmt"
	"strings"
	"time"
)

type Operation int

const (
	In Operation = iota
	LessThan
	GreaterThan
)

type Filter struct {
	Key      string
	Values   []interface{}
	Operator Operation
}

func (f Filter) String() string {
	switch f.Operator {
	case LessThan:
		return f.stringifyOperator("<")
	case GreaterThan:
		return f.stringifyOperator(">")
	case In:
		fallthrough
	default:
		return f.stringifyIn()
	}
}

func (f Filter) stringifyIn() string {
	quoted := make([]string, len(f.Values))
	for id, v := range f.Values {
		// https://stackoverflow.com/questions/37782278/fully-parsing-timestamps-in-golang
		t, ok := v.(time.Time)
		if ok {
			quoted[id] = fmt.Sprintf("'%v'", t.Format(time.RFC3339))
			continue
		}

		quoted[id] = fmt.Sprintf("'%v'", v)
	}

	return fmt.Sprintf("%s in (%s)", f.Key, strings.Join(quoted, ","))
}

func (f Filter) stringifyOperator(op string) string {
	out := ""

	for id, filter := range f.Values {
		if id > 0 {
			out += " and "
		}

		t, ok := filter.(time.Time)
		if ok {
			out += fmt.Sprintf("%s %s '%v'", f.Key, op, t.Format(time.RFC3339))

			continue
		}

		out += fmt.Sprintf("%s %s '%v'", f.Key, op, filter)
	}

	return out
}
