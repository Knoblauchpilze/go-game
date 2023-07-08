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

	logger.Infof("connection to %s closed", db.config)

	return nil
}

func (db *postgresDb) Query(query Query) Rows {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.pool == nil {
		return newRows(nil, errors.NewCode(errors.ErrDbConnectionInvalid))
	}

	if !query.Valid() {
		return newRows(nil, errors.NewCode(errors.ErrInvalidQuery))
	}

	sqlQuery := query.ToSql()

	var err error
	rows, err := db.pool.Query(sqlQuery)
	if err != nil {
		return newRows(nil, errors.WrapCode(err, errors.ErrDbRequestFailed))
	}

	return newRows(rows, nil)
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

	var err error
	out.tag, err = db.pool.Exec(sqlQuery)
	if err != nil {
		out.Err = errors.WrapCode(err, errors.ErrDbRequestFailed)
	}

	return out
}
