package db

import (
	"sync"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/jackc/pgx"
)

const dbHost = "localhost"
const dbPort = uint16(5500)
const dbName = "user_service_db"
const dbUser = "user_service_administrator"
const dbPassword = "Ww76hQWbbt7zi2ItM6cNo4YYT"

const dbConnectionsPoolSize = 1
const dbHealthcheckInterval = 10 * time.Second

type dbRepository struct {
	pool *pgx.ConnPool
	lock sync.Mutex
}

func NewRepository() Repository {
	// Create the DB object.
	repo := dbRepository{
		nil,
		sync.Mutex{},
	}

	repo.createPoolAttempt()

	ticker := time.NewTicker(dbHealthcheckInterval)
	go func() {
		for range ticker.C {
			repo.healthcheck()
		}
	}()

	return &repo
}

func (repo *dbRepository) Foo() {
	logger.Infof("salut")
}

func (repo *dbRepository) createPoolAttempt() bool {
	logger.Infof("Connection attempt to %v at %v:%v, user %v", dbName, dbHost, dbPort, dbUser)

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
		logger.Warnf("Failed to connect to %v at %v:%v (err: %v)", dbName, dbHost, dbPort, err)
		return false
	}

	logger.Infof("Connected to %v at %v:%v", dbName, dbHost, dbPort)

	repo.lock.Lock()
	func() {
		defer repo.lock.Unlock()
		repo.pool = pool
	}()

	return true
}

func (repo *dbRepository) healthcheck() {
	dbIsNil := false
	var stat pgx.ConnPoolStat

	repo.lock.Lock()
	func() {
		defer repo.lock.Unlock()

		dbIsNil = (repo.pool == nil)
		if !dbIsNil {
			stat = repo.pool.Stat()
		}
	}()

	if dbIsNil || stat.CurrentConnections == 0 {
		repo.createPoolAttempt()
	}
}
