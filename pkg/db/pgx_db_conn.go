package db

import "github.com/jackc/pgx"

type pgxDbConn interface {
	Close()
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error)
}
