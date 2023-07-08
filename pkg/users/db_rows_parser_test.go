package users

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRowParser_ScanRow(t *testing.T) {
	assert := assert.New(t)

	m := &mockScannable{}
	urp := userRowParser{}

	err := urp.ScanRow(m)
	assert.Nil(err)
	assert.Equal(1, m.scanCalled)
}

func TestUserRowParser_ScanRow_Error(t *testing.T) {
	assert := assert.New(t)

	m := &mockScannable{
		scanErr: fmt.Errorf("someError"),
	}
	urp := userRowParser{}

	err := urp.ScanRow(m)
	assert.Contains(err.Error(), "someError")
	assert.Equal(1, m.scanCalled)
}

func TestUserIdsParser_ScanRow(t *testing.T) {
	assert := assert.New(t)

	m := &mockScannable{}
	urp := userIdsParser{}

	err := urp.ScanRow(m)
	assert.Nil(err)
	assert.Equal(1, m.scanCalled)
}

func TestUserIdsParser_ScanRow_Error(t *testing.T) {
	assert := assert.New(t)

	m := &mockScannable{
		scanErr: fmt.Errorf("someError"),
	}
	urp := userIdsParser{}

	err := urp.ScanRow(m)
	assert.Contains(err.Error(), "someError")
	assert.Equal(1, m.scanCalled)
}

type mockScannable struct {
	scanCalled int
	scanErr    error
}

func (m *mockScannable) Scan(dest ...interface{}) error {
	m.scanCalled++
	return m.scanErr
}
