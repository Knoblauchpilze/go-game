package db

import (
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestExtractRows(t *testing.T) {
	assert := assert.New(t)

	_, err := extractRows("not-a-number")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))

	n, err := extractRows("32")
	assert.Equal(32, n)
	assert.Nil(err)
}

func TestExtractAffectedRowsFromCommandTag_Invalid(t *testing.T) {
	assert := assert.New(t)

	tag := pgx.CommandTag("invalidTag")
	_, err := extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))

	tag = pgx.CommandTag("invalidTag with many fields")
	_, err = extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))
}

func TestExtractAffectedRowsFromCommandTag_UnknownTag(t *testing.T) {
	assert := assert.New(t)

	tag := pgx.CommandTag("UNKNOWN_VERB ok fields")
	_, err := extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrUnknownSqlCommandTag))
}

func TestExtractAffectedRowsFromCommandTag_Insert(t *testing.T) {
	assert := assert.New(t)

	tag := pgx.CommandTag("INSERT 24")
	_, err := extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))

	tag = pgx.CommandTag("INSERT 0 not-a-number")
	_, err = extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))

	tag = pgx.CommandTag("INSERT 0 37")
	n, err := extractAffectedRowsFromCommandTag(tag)
	assert.Nil(err)
	assert.Equal(37, n)
}

func TestExtractAffectedRowsFromCommandTag_Deleted(t *testing.T) {
	assert := assert.New(t)

	tag := pgx.CommandTag("DELETE not-a-number")
	_, err := extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))

	tag = pgx.CommandTag("DELETE 23 additional-argument")
	_, err = extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))

	tag = pgx.CommandTag("DELETE 24")
	n, err := extractAffectedRowsFromCommandTag(tag)
	assert.Nil(err)
	assert.Equal(24, n)
}
