package db

import (
	"sync"

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

func (db *postgresDb) Connect() error {
	logger.Infof("connection attempt to %s", db.config)

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

	pool, err := db.config.creationFunc(pgxConf)
	if err != nil {
		return errors.WrapCode(err, errors.ErrDbConnectionFailed)
	}

	logger.Infof("connected to %s", db.config)

	db.lock.Lock()
	func() {
		defer db.lock.Unlock()
		db.pool = pool
	}()

	return nil
}

func (db *postgresDb) Disconnect() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.pool == nil {
		return nil
	}

	db.pool.Close()
	db.pool = nil

	return nil
}

func (db *postgresDb) Query(query Query) QueryRows {
	db.lock.Lock()
	defer db.lock.Unlock()

	out := &queryRowsImpl{}

	if db.pool == nil {
		out.err = errors.NewCode(errors.ErrDbConnectionInvalid)
		return out
	}

	if !query.Valid() {
		out.err = errors.NewCode(errors.ErrInvalidQuery)
		return out
	}

	sqlQuery := query.ToSql()

	if query.Verbose() {
		logger.Tracef("executing query: %s", sqlQuery)
	}

	out.rows, out.err = db.pool.Query(sqlQuery)

	return out
}

func (db *postgresDb) Execute(query Query) Result {
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
		logger.Tracef("executing script: %s", sqlQuery)
	}

	out.tag, out.Err = db.pool.Exec(sqlQuery)

	return out
}
