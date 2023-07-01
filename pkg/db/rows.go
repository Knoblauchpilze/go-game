package db

import "github.com/jackc/pgx"

type Rows struct {
	rows *pgx.Rows
	Err  error
}

func (r Rows) Next() bool {
	return r.rows.Next()
}

func (r Rows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r Rows) Close() {
	if r.rows != nil {
		r.rows.Close()
	}
}
