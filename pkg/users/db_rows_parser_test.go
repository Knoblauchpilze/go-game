package users

import (
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
		scanErr: errDefault,
	}
	urp := userRowParser{}

	err := urp.ScanRow(m)
	assert.Equal(errDefault, err)
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
		scanErr: errDefault,
	}
	urp := userIdsParser{}

	err := urp.ScanRow(m)
	assert.Equal(errDefault, err)
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
