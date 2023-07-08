package db

import (
	"github.com/KnoblauchPilze/go-game/pkg/common"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type Rows interface {
	Err() error
	Close()
	Empty() bool

	GetSingleValue(parser Parser) error
	GetAll(parser Parser) error
}

type Scannable interface {
	Scan(dest ...interface{}) error
}

type Parser interface {
	Parse(row Scannable) error
}

type sqlRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
}

type rowsImpl struct {
	rows sqlRows
	next bool
	err  error
}

func newRows(rows sqlRows, err error) Rows {
	r := rowsImpl{
		rows: rows,
		err:  err,
	}

	if !common.IsInterfaceNil(r.rows) && r.err == nil {
		r.next = r.rows.Next()
	}

	return &r
}

func (r *rowsImpl) Err() error {
	return r.err
}

func (r *rowsImpl) Close() {
	if r.rows != nil {
		r.rows.Close()
	}
}

func (r *rowsImpl) Empty() bool {
	return r.rows == nil || !r.next
}

func (r *rowsImpl) GetSingleValue(parser Parser) error {
	if err := r.assertValidStateOrReturnError(); err != nil {
		return err
	}

	defer r.Close()

	if err := parser.Parse(r.rows); err != nil {
		return errors.WrapCode(err, errors.ErrSqlRowParsingFailed)
	}

	r.next = r.rows.Next()
	if r.next {
		return errors.NewCode(errors.ErrMultiValuedDbElement)
	}

	return nil
}

func (r *rowsImpl) GetAll(parser Parser) error {
	if err := r.assertValidStateOrReturnError(); err != nil {
		return err
	}

	defer r.Close()

	for r.next {
		if err := parser.Parse(r.rows); err != nil {
			return errors.WrapCode(err, errors.ErrSqlRowParsingFailed)
		}

		r.next = r.rows.Next()
	}

	return nil
}

func (r *rowsImpl) assertValidStateOrReturnError() error {
	if err := r.Err(); err != nil {
		return err
	}
	if r.Empty() {
		return errors.NewCode(errors.ErrNoRowsReturnedForSqlQuery)
	}

	return nil
}
