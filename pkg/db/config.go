package db

import (
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

type CreationFunc func(config pgx.ConnPoolConfig) (pgxDbFacade, error)

type Config struct {
	DbHost                string
	DbPort                uint16
	DbName                string
	DbUser                string
	DbPassword            string
	DbConnectionsPoolSize uint
	DbConnectionTimeout   time.Duration
	DbQueryTimeout        time.Duration
	creationFunc          CreationFunc
}

func NewConfig() Config {
	conf := Config{
		creationFunc: func(config pgx.ConnPoolConfig) (pgxDbFacade, error) {
			return pgx.NewConnPool(config)
		},
	}

	return conf
}

func (c Config) Valid() bool {
	if len(c.DbHost) == 0 {
		return false
	}
	if len(c.DbName) == 0 {
		return false
	}
	if len(c.DbUser) == 0 {
		return false
	}
	if len(c.DbPassword) == 0 {
		return false
	}
	if c.DbConnectionsPoolSize < 2 {
		return false
	}

	return true
}

func (c Config) String() string {
	return fmt.Sprintf("%v %v@%v:%v", c.DbName, c.DbUser, c.DbHost, c.DbPort)
}
