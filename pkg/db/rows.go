package db

import "github.com/jackc/pgx"

type QueryRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
	Close()
}

type queryRowsImpl struct {
	rows *pgx.Rows
	err  error
}

func (qr *queryRowsImpl) Next() bool {
	return qr.rows.Next()
}

func (qr *queryRowsImpl) Scan(dest ...interface{}) error {
	return qr.rows.Scan(dest...)
}

func (qr *queryRowsImpl) Err() error {
	return qr.err
}

func (qr queryRowsImpl) Close() {
	if qr.rows != nil {
		qr.rows.Close()
	}
}
