package db

import (
	"fmt"
	"strings"
)

type selectQuery struct {
	props    []string
	with     Query
	withName string
	table    string
	filters  []Filter
	ordering string
	verbose  bool
}

func (q selectQuery) Valid() bool {
	return len(q.props) > 0 && len(q.table) > 0
}

func (q selectQuery) Verbose() bool {
	return q.verbose
}

func (q selectQuery) ToSql() (string, error) {
	with := ""
	if q.with != nil {
		withClause, err := q.with.ToSql()
		if err != nil {
			return "", err
		}
		with = fmt.Sprintf("with %s as (%s)", q.withName, withClause)
	}

	str := fmt.Sprintf("%s select %s from %s", with, strings.Join(q.props, ", "), q.table)

	if len(q.filters) > 0 {
		str += " where"

		for id, filter := range q.filters {
			if id > 0 {
				str += " and"
			}

			str += fmt.Sprintf(" %s", filter)
		}
	}

	str += fmt.Sprintf(" %s", q.ordering)

	return str, nil
}
