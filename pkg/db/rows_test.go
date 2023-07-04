package db

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mockRows struct {
	count        int
	numberOfRows int
	scanError    error
	closeCalls   int
}

func (m *mockRows) Next() bool {
	out := m.count < m.numberOfRows
	m.count++
	return out
}

func (m *mockRows) Scan(dest ...interface{}) error {
	return m.scanError
}

func (m *mockRows) Close() {
	m.closeCalls++
}

func TestRows_Err(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, nil)
	assert.Nil(r.Err())

	r = newRows(nil, fmt.Errorf("someError"))
	assert.Equal("someError", r.Err().Error())
}

func TestRows_Close(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, nil)
	assert.NotPanics(r.Close)

	var m mockRows
	r = newRows(&m, nil)
	r.Close()
	assert.Equal(1, m.closeCalls)
}

func TestRows_Empty(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, nil)
	assert.True(r.Empty())

	m := &mockRows{}
	r = newRows(m, nil)
	assert.True(r.Empty())

	m = &mockRows{
		numberOfRows: 1,
	}
	r = newRows(m, nil)
	assert.False(r.Empty())
}

func TestRows_GetSingleValue_InvalidPreconditions(t *testing.T) {
	assert := assert.New(t)

	calls := 0
	scan := func(row Scannable) error {
		calls++
		return nil
	}

	r := newRows(nil, fmt.Errorf("someError"))
	err := r.GetSingleValue(scan)
	assert.Equal("someError", err.Error())
	assert.Equal(0, calls)

	r = newRows(nil, nil)
	err = r.GetSingleValue(scan)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoRowsReturnedForSqlQuery))
	assert.Equal(0, calls)
}

func TestRows_GetSingleValue(t *testing.T) {
	assert := assert.New(t)

	calls := 0
	scan := func(row Scannable) error {
		calls++
		return nil
	}

	m := &mockRows{
		numberOfRows: 1,
	}
	r := newRows(m, nil)
	err := r.GetSingleValue(scan)
	assert.Nil(err)
	assert.Equal(1, calls)
}

func TestRows_GetSingleValue_ScannerError(t *testing.T) {
	assert := assert.New(t)

	calls := 0
	scan := func(row Scannable) error {
		calls++
		return fmt.Errorf("someError")
	}

	m := &mockRows{
		numberOfRows: 1,
	}
	r := newRows(m, nil)
	err := r.GetSingleValue(scan)
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlRowParsingFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
	assert.Equal(1, calls)
}

func TestRows_GetSingleValue_WithMultipleValues(t *testing.T) {
	assert := assert.New(t)

	calls := 0
	scan := func(row Scannable) error {
		calls++
		return nil
	}

	m := &mockRows{
		numberOfRows: 2,
	}
	r := newRows(m, nil)
	err := r.GetSingleValue(scan)
	assert.True(errors.IsErrorWithCode(err, errors.ErrMultiValuedDbElement))
}

func TestRows_GetAll_InvalidPreconditions(t *testing.T) {
	assert := assert.New(t)

	calls := 0
	scan := func(row Scannable) error {
		calls++
		return nil
	}

	r := newRows(nil, fmt.Errorf("someError"))
	err := r.GetAll(scan)
	assert.Equal("someError", err.Error())
	assert.Equal(0, calls)

	r = newRows(nil, nil)
	err = r.GetAll(scan)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoRowsReturnedForSqlQuery))
	assert.Equal(0, calls)
}

func TestRows_GetAll(t *testing.T) {
	assert := assert.New(t)

	calls := 0
	scan := func(row Scannable) error {
		calls++
		return nil
	}

	m := &mockRows{
		numberOfRows: 2,
	}
	r := newRows(m, nil)
	err := r.GetAll(scan)
	assert.Nil(err)
	assert.Equal(2, calls)
}

func TestRows_GetAll_ScannerError(t *testing.T) {
	assert := assert.New(t)

	calls := 0
	scan := func(row Scannable) error {
		calls++
		return fmt.Errorf("someError")
	}

	m := &mockRows{
		numberOfRows: 2,
	}
	r := newRows(m, nil)
	err := r.GetAll(scan)
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlRowParsingFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
	assert.Equal(1, calls)
}

// type Rows interface {

// 	GetAll(outSlice interface{}) error
// }