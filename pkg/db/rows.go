package db

import "github.com/jackc/pgx"

// https://github.com/jackc/pgx/issues/878
type QueryRows interface {
	Err() error
	Close()
	Empty() bool

	//Traverse() error

	Next() bool
	Scan(dest ...interface{}) error
}

type queryRowsImpl struct {
	rows *pgx.Rows
	err  error
}

func (qr queryRowsImpl) Err() error {
	return qr.err
}

func (qr queryRowsImpl) Close() {
	if qr.rows != nil {
		qr.rows.Close()
	}
}

func (qr queryRowsImpl) Empty() bool {
	return qr.rows != nil && qr.rows.Next()
}

func (qr queryRowsImpl) Next() bool {
	return qr.rows != nil && qr.rows.Next()
}

func (qr queryRowsImpl) Scan(dest ...interface{}) error {
	return qr.rows.Scan(dest...)
}
