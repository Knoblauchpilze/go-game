package db

import (
	"net"
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestConfig_New(t *testing.T) {
	assert := assert.New(t)

	conf := NewConfig()
	assert.NotNil(conf.creationFunc)
}

func TestConfig_Create(t *testing.T) {
	assert := assert.New(t)

	conf := NewConfig()

	pgxConf := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host: "host",
			Dial: func(network, addr string) (net.Conn, error) {
				return nil, errDefault
			},
		},
	}

	_, err := conf.creationFunc(pgxConf)
	assert.Equal(errDefault, err)
}

func TestConfig_Valid(t *testing.T) {
	assert := assert.New(t)

	conf := Config{}
	assert.False(conf.Valid())

	conf.DbHost = "host"
	assert.False(conf.Valid())

	conf.DbPort = 32
	assert.False(conf.Valid())

	conf.DbName = "database"
	assert.False(conf.Valid())

	conf.DbUser = "user"
	assert.False(conf.Valid())

	conf.DbPassword = "password"
	assert.False(conf.Valid())

	conf.DbConnectionsPoolSize = 2
	assert.True(conf.Valid())
}

func TestConfig_String(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		DbHost:                "host",
		DbPort:                32,
		DbName:                "database",
		DbUser:                "user",
		DbPassword:            "password",
		DbConnectionsPoolSize: 2,
	}

	assert.Equal("database user@host:32", conf.String())
}
