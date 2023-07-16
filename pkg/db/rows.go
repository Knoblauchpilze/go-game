package db

import (
	"github.com/KnoblauchPilze/go-game/pkg/common"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type Rows interface {
	Err() error
	Close()
	Empty() bool

	GetSingleValue(parser RowParser) error
	GetAll(parser RowParser) error
}

type Scannable interface {
	Scan(dest ...interface{}) error
}

type RowParser interface {
	ScanRow(row Scannable) error
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

func (r *rowsImpl) GetSingleValue(parser RowParser) error {
	if err := r.Err(); err != nil {
		return err
	}
	if r.Empty() {
		return errors.NewCode(errors.ErrNoRowsReturnedForSqlQuery)
	}

	defer r.Close()

	if err := parser.ScanRow(r.rows); err != nil {
		return errors.WrapCode(err, errors.ErrSqlRowParsingFailed)
	}

	r.next = r.rows.Next()
	if r.next {
		return errors.NewCode(errors.ErrMultiValuedDbElement)
	}

	return nil
}

func (r *rowsImpl) GetAll(parser RowParser) error {
	if err := r.Err(); err != nil {
		return err
	}

	defer r.Close()

	for r.next {
		if err := parser.ScanRow(r.rows); err != nil {
			return errors.WrapCode(err, errors.ErrSqlRowParsingFailed)
		}

		r.next = r.rows.Next()
	}

	return nil
}
