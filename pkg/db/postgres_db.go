package db

import (
	"sync"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/jackc/pgx"
)

const dbHost = "localhost"
const dbPort = uint16(5500)
const dbName = "user_service_db"
const dbUser = "user_service_administrator"
const dbPassword = "Ww76hQWbbt7zi2ItM6cNo4YYT"

const dbConnectionsPoolSize = 1
const dbHealthcheckInterval = 5 * time.Second

// https://betterprogramming.pub/how-to-work-with-sql-in-go-ca8bc0b30722
// https://www.sohamkamani.com/golang/sql-database/
type postgresDb struct {
	pool *pgx.ConnPool
	lock sync.Mutex
}

func NewPostgresDatabase() Database {
	db := postgresDb{
		nil,
		sync.Mutex{},
	}

	db.createPoolAttempt()

	ticker := time.NewTicker(dbHealthcheckInterval)
	go func() {
		for range ticker.C {
			db.healthcheck()
		}
	}()

	return &db
}

func (db *postgresDb) Query(query Query) (Rows, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.pool == nil {
		return Rows{}, errors.NewCode(errors.ErrDbConnectionInvalid)
	}

	if !query.Valid() {
		return Rows{}, errors.NewCode(errors.ErrInvalidQuery)
	}

	sqlQuery, err := query.ToSql()
	if err != nil {
		return Rows{}, errors.NewCode(errors.ErrSqlTranslation)
	}

	if query.Verbose() {
		logger.Tracef("executing query: %s", sqlQuery)
	}

	var out Rows
	out.rows, out.Err = db.pool.Query(sqlQuery)

	return out, nil
}

func (db *postgresDb) Execute(query Query) (Result, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.pool == nil {
		return Result{}, errors.NewCode(errors.ErrDbConnectionInvalid)
	}

	if !query.Valid() {
		return Result{}, errors.NewCode(errors.ErrInvalidQuery)
	}

	sqlQuery, err := query.ToSql()
	if err != nil {
		return Result{}, errors.NewCode(errors.ErrSqlTranslation)
	}

	if query.Verbose() {
		logger.Tracef("executing script: %s", sqlQuery)
	}

	var out Result
	out.tag, out.Err = db.pool.Exec(sqlQuery)

	return out, nil
}

func (db *postgresDb) createPoolAttempt() bool {
	logger.Infof("connection attempt to %v at %v:%v, user %v", dbName, dbHost, dbPort, dbUser)

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     dbHost,
			Database: dbName,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
		},
		MaxConnections: dbConnectionsPoolSize,
		AcquireTimeout: 0,
	})

	if err != nil {
		logger.Warnf("failed to connect to %v at %v:%v (err: %v)", dbName, dbHost, dbPort, err)
		return false
	}

	logger.Infof("connected to %v at %v:%v", dbName, dbHost, dbPort)

	db.lock.Lock()
	func() {
		defer db.lock.Unlock()
		db.pool = pool
	}()

	return true
}

func (db *postgresDb) healthcheck() {
	dbIsNil := false
	var stat pgx.ConnPoolStat

	db.lock.Lock()
	func() {
		defer db.lock.Unlock()

		dbIsNil = (db.pool == nil)
		if !dbIsNil {
			stat = db.pool.Stat()
		}
	}()

	logger.Debugf("stats: %+v", stat)

	if dbIsNil || stat.CurrentConnections == 0 {
		db.createPoolAttempt()
	}
}
