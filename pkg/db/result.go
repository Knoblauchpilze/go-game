package db

import "github.com/jackc/pgx"

type Result struct {
	tag pgx.CommandTag
	Err error
}
