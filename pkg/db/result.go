package db

import (
	"github.com/jackc/pgx"
)

type Result interface {
	Err() error
	AffectedRows() int
}

type resultImpl struct {
	affectedRows int
	err          error
}

func newResult(tag pgx.CommandTag, sqlErr error) Result {
	r := resultImpl{
		err: sqlErr,
	}

	var err error
	r.affectedRows, err = extractAffectedRowsFromCommandTag(tag)
	if err != nil && r.err == nil {
		r.err = err
	}

	return &r
}

func (r *resultImpl) Err() error {
	return r.err
}

func (r *resultImpl) AffectedRows() int {
	return r.affectedRows
}
