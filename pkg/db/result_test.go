package db

import (
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestResult_New(t *testing.T) {
	assert := assert.New(t)

	out := newResult("INSERT 0 12", nil)
	assert.Nil(out.Err())
	assert.Equal(12, out.AffectedRows())
}

func TestResult_New_WithError(t *testing.T) {
	assert := assert.New(t)

	out := newResult("INSERT 0 12", errDefault)
	assert.Equal(errDefault, out.Err())
	assert.Equal(12, out.AffectedRows())
}

func TestResult_New_InvalidTagAndError(t *testing.T) {
	assert := assert.New(t)

	out := newResult("not-a-valid-tag", errDefault)
	assert.Equal(errDefault, out.Err())
	assert.Equal(0, out.AffectedRows())
}

func TestResult_New_InvalidTag(t *testing.T) {
	assert := assert.New(t)

	out := newResult("not-a-valid-tag", nil)
	assert.True(errors.IsErrorWithCode(out.Err(), errors.ErrInvalidSqlCommandTag))
	assert.Equal(0, out.AffectedRows())
}
