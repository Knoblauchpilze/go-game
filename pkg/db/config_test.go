package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresDatabase_New(t *testing.T) {
	assert := assert.New(t)

	conf := NewConfig()
	assert.NotNil(conf.creationFunc)
}

func TestPostgresDatabase_Valid(t *testing.T) {
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

func TestPostgresDatabase_String(t *testing.T) {
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
