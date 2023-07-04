package db

import (
	"github.com/KnoblauchPilze/go-game/pkg/common"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type Scannable interface {
	Scan(dest ...interface{}) error
}

type ScanRow func(row Scannable) error

// https://github.com/jackc/pgx/issues/878
type Rows interface {
	Err() error
	Close()
	Empty() bool

	GetSingleValue(scanner ScanRow) error
	GetAll(scanner ScanRow) error
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

func (r *rowsImpl) GetSingleValue(scanner ScanRow) error {
	if err := r.assertValidStateOrReturnError(); err != nil {
		return err
	}

	if err := scanner(r.rows); err != nil {
		return errors.WrapCode(err, errors.ErrSqlRowParsingFailed)
	}

	r.next = r.rows.Next()
	if r.next {
		return errors.NewCode(errors.ErrMultiValuedDbElement)
	}

	return nil
}

func (r *rowsImpl) GetAll(scanner ScanRow) error {
	if err := r.assertValidStateOrReturnError(); err != nil {
		return err
	}

	for r.next {
		if err := scanner(r.rows); err != nil {
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

// func (qr queryRowsImpl) Next() bool {
// 	return qr.rows != nil && qr.rows.Next()
// }

// func (qr queryRowsImpl) Scan(dest ...interface{}) error {
// 	return qr.rows.Scan(dest...)
// }
