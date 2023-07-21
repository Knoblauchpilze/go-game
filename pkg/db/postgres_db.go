package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/KnoblauchPilze/go-game/pkg/common"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/jackc/pgx"
)

// https://www.sohamkamani.com/golang/sql-database/
// https://betterprogramming.pub/how-to-work-with-sql-in-go-ca8bc0b30722
type postgresDb struct {
	config Config
	pool   pgxDbFacade
	lock   sync.Mutex
}

func NewPostgresDatabase(conf Config) Database {
	db := postgresDb{
		config: conf,
	}

	return &db
}

func (db *postgresDb) Connect(ctx context.Context) error {
	logger.ScopedInfof(ctx, "connection attempt to %s", db.config)

	fmt.Printf("%+v\n", db.config.DbConnectionTimeout)

	pgxConf := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     db.config.DbHost,
			Database: db.config.DbName,
			Port:     db.config.DbPort,
			User:     db.config.DbUser,
			Password: db.config.DbPassword,
		},
		MaxConnections: int(db.config.DbConnectionsPoolSize),
		AcquireTimeout: 0,
	}

	var pool pgxDbFacade
	connFunc := func() error {
		var err error
		pool, err = db.config.creationFunc(pgxConf)
		return err
	}

	res := common.ExecuteWithContext(connFunc, ctx, db.config.DbConnectionTimeout)
	if err := res.Error(); err != nil {
		if err == context.DeadlineExceeded {
			return errors.WrapCode(err, errors.ErrDbConnectionTimeout)
		}
		return errors.WrapCode(err, errors.ErrDbConnectionFailed)
	}

	logger.ScopedInfof(ctx, "connected to %s", db.config)

	db.lock.Lock()
	func() {
		defer db.lock.Unlock()
		db.pool = pool
	}()

	return nil
}

func (db *postgresDb) Disconnect(ctx context.Context) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.pool == nil {
		return nil
	}

	db.pool.Close()
	db.pool = nil

	logger.ScopedInfof(ctx, "connection to %s closed", db.config)

	return nil
}

func (db *postgresDb) Query(ctx context.Context, query Query) Rows {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.pool == nil {
		return newRows(nil, errors.NewCode(errors.ErrDbConnectionInvalid))
	}

	if !query.Valid() {
		return newRows(nil, errors.NewCode(errors.ErrInvalidQuery))
	}

	sqlQuery := query.ToSql()
	if query.Verbose() {
		logger.ScopedTracef(ctx, "executing: %s", query.ToSql())
	}

	var err error
	rows, err := db.pool.Query(sqlQuery)
	if err != nil {
		return newRows(nil, errors.WrapCode(err, errors.ErrDbRequestFailed))
	}

	return newRows(rows, nil)
}

func (db *postgresDb) Execute(ctx context.Context, query Query) Result {
	db.lock.Lock()
	defer db.lock.Unlock()

	var out Result

	if db.pool == nil {
		out.Err = errors.NewCode(errors.ErrDbConnectionInvalid)
		return out
	}

	if !query.Valid() {
		out.Err = errors.NewCode(errors.ErrInvalidQuery)
		return out
	}

	sqlQuery := query.ToSql()
	if query.Verbose() {
		logger.ScopedTracef(ctx, "executing: %s", query.ToSql())
	}

	var err error
	out.tag, err = db.pool.Exec(sqlQuery)
	if err != nil {
		out.Err = errors.WrapCode(err, errors.ErrDbRequestFailed)
	}

	return out
}
