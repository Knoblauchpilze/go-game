package db

import (
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestExtractAffectedRowsFromCommandTag_Invalid(t *testing.T) {
	assert := assert.New(t)

	tag := pgx.CommandTag("invalidTag")
	_, err := extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))

	tag = pgx.CommandTag("insert 0 not-a-number")
	_, err = extractAffectedRowsFromCommandTag(tag)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidSqlCommandTag))
	cause := errors.Unwrap(err)
	assert.Contains(cause.Error(), "parsing \"not-a-number\": invalid syntax")
}

func TestExtractAffectedRowsFromCommandTag(t *testing.T) {
	assert := assert.New(t)

	tag := pgx.CommandTag("insert 12 24")
	n, err := extractAffectedRowsFromCommandTag(tag)
	assert.Nil(err)
	assert.Equal(24, n)
}
