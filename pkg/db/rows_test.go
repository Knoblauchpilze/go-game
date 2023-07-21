package db

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

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

	var m mockSqlRows
	r = newRows(&m, nil)
	r.Close()
	assert.Equal(int32(1), m.closeCalls.Load())
}

func TestRows_Empty(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, nil)
	assert.True(r.Empty())

	m := &mockSqlRows{}
	r = newRows(m, nil)
	assert.True(r.Empty())

	m = &mockSqlRows{
		numberOfRows: 1,
	}
	r = newRows(m, nil)
	assert.False(r.Empty())
}

func TestRows_GetSingleValue_InvalidPreconditions(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}

	r := newRows(nil, fmt.Errorf("someError"))
	err := r.GetSingleValue(mp)
	assert.Equal("someError", err.Error())
	assert.Equal(0, mp.parseCalled)

	r = newRows(nil, nil)
	err = r.GetSingleValue(mp)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoRowsReturnedForSqlQuery))
	assert.Equal(0, mp.parseCalled)
}

func TestRows_GetSingleValue(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}
	m := &mockSqlRows{
		numberOfRows: 1,
	}
	r := newRows(m, nil)
	err := r.GetSingleValue(mp)
	assert.Nil(err)
	assert.Equal(1, mp.parseCalled)
}

func TestRows_GetSingleValue_CallsClose(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}
	m := &mockSqlRows{
		numberOfRows: 1,
	}
	r := newRows(m, nil)
	r.GetSingleValue(mp)
	assert.Equal(int32(1), m.closeCalls.Load())
}

func TestRows_GetSingleValue_ScannerError(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{
		parseErr: fmt.Errorf("someError"),
	}

	m := &mockSqlRows{
		numberOfRows: 1,
	}
	r := newRows(m, nil)
	err := r.GetSingleValue(mp)
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlRowParsingFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
	assert.Equal(1, mp.parseCalled)
}

func TestRows_GetSingleValue_WithMultipleValues(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}
	m := &mockSqlRows{
		numberOfRows: 2,
	}
	r := newRows(m, nil)
	err := r.GetSingleValue(mp)
	assert.True(errors.IsErrorWithCode(err, errors.ErrMultiValuedDbElement))
}

func TestRows_GetAll_InvalidPreconditions(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}
	r := newRows(nil, fmt.Errorf("someError"))
	err := r.GetAll(mp)
	assert.Equal("someError", err.Error())
	assert.Equal(0, mp.parseCalled)
}

func TestRows_GetAll(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}
	m := &mockSqlRows{
		numberOfRows: 2,
	}
	r := newRows(m, nil)
	err := r.GetAll(mp)
	assert.Nil(err)
	assert.Equal(2, mp.parseCalled)
}

func TestRows_GetAll_NoData(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}
	m := &mockSqlRows{
		numberOfRows: 0,
	}
	r := newRows(m, nil)
	err := r.GetAll(mp)
	assert.Nil(err)
	assert.Equal(0, mp.parseCalled)
}

func TestRows_GetAll_CallsClose(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}
	m := &mockSqlRows{
		numberOfRows: 0,
	}
	r := newRows(m, nil)
	r.GetAll(mp)
	assert.Equal(int32(1), m.closeCalls.Load())

	m.closeCalls.Store(0)
	m.numberOfRows = 2
	r = newRows(m, nil)
	r.GetAll(mp)
	assert.Equal(int32(1), m.closeCalls.Load())
}

func TestRows_GetAll_ScannerError(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{parseErr: fmt.Errorf("someError")}
	m := &mockSqlRows{
		numberOfRows: 2,
	}
	r := newRows(m, nil)
	err := r.GetAll(mp)
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlRowParsingFailed))
	cause := errors.Unwrap(err)
	assert.Equal("someError", cause.Error())
	assert.Equal(1, mp.parseCalled)
}

type mockSqlRows struct {
	count        int
	numberOfRows int
	scanError    error
	closeCalls   atomic.Int32
}

func (m *mockSqlRows) Next() bool {
	out := m.count < m.numberOfRows
	m.count++
	return out
}

func (m *mockSqlRows) Scan(dest ...interface{}) error {
	return m.scanError
}

func (m *mockSqlRows) Close() {
	m.closeCalls.Add(1)
}

type mockParser struct {
	parseCalled int
	parseErr    error
}

func (m *mockParser) ScanRow(row Scannable) error {
	m.parseCalled++
	return m.parseErr
}
