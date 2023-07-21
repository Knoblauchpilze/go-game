package db

import "github.com/jackc/pgx"

type pgxDbFacade interface {
	Close()
	Query(sql string, args ...interface{}) (sqlRows, error)
	Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error)
}

type pgxDbFacadeImpl struct {
	pool *pgx.ConnPool
}

func newPgxDbFacadeImpl(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
	pool, err := pgx.NewConnPool(config)
	f := pgxDbFacadeImpl{
		pool: pool,
	}
	return &f, err
}

func (f *pgxDbFacadeImpl) Close() {
	f.pool.Close()
}

func (f *pgxDbFacadeImpl) Query(sql string, args ...interface{}) (sqlRows, error) {

	return f.pool.Query(sql, args...)
}

func (f *pgxDbFacadeImpl) Exec(sql string, args ...interface{}) (pgx.CommandTag, error) {
	return f.pool.Exec(sql, args...)
}
