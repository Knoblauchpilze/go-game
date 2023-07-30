package db

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestQueryExecutor_runQueryAndReturnRows_BuildError(t *testing.T) {
	assert := assert.New(t)

	mdb := &mockDb{}
	mqb := mockQueryBuilder{
		buildErr: errDefault,
	}

	qe := queryExecutorImpl{
		db: mdb,
	}

	_, err := qe.runQueryAndReturnRows(context.TODO(), mqb)
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal(errDefault, cause)
}

func TestQueryExecutor_runQueryAndReturnRows_QueryError(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mdb := &mockDb{
		rows: &mockRows{
			err: errDefault,
		},
	}

	qe := queryExecutorImpl{
		db: mdb,
	}

	_, err := qe.runQueryAndReturnRows(context.TODO(), mqb)
	assert.Equal(errDefault, err)
}

func TestQueryExecutor_runQueryAndReturnRows(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockRows{}
	mdb := &mockDb{
		rows: mr,
	}

	qe := queryExecutorImpl{
		db: mdb,
	}

	_, err := qe.runQueryAndReturnRows(context.TODO(), mqb)
	assert.Nil(err)
	assert.Equal(1, mdb.queryCalls)
}

func TestQueryExecutor_executeQueryAndReturn_BuildError(t *testing.T) {
	assert := assert.New(t)

	mdb := &mockDb{}
	mqb := mockQueryBuilder{
		buildErr: errDefault,
	}

	qe := queryExecutorImpl{
		db: mdb,
	}

	_, err := qe.executeQueryAndReturn(context.TODO(), mqb)
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbRequestCreationFailed))
	cause := errors.Unwrap(err)
	assert.Equal(errDefault, cause)
}

func TestQueryExecutor_executeQueryAndReturn_QueryError(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mdb := &mockDb{
		result: &mockResult{
			err: errDefault,
		},
	}

	qe := queryExecutorImpl{
		db: mdb,
	}

	_, err := qe.executeQueryAndReturn(context.TODO(), mqb)
	assert.Equal(errDefault, err)
}

func TestQueryExecutor_executeQueryAndReturn(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockResult{}
	mdb := &mockDb{
		result: mr,
	}

	qe := queryExecutorImpl{
		db: mdb,
	}

	_, err := qe.executeQueryAndReturn(context.TODO(), mqb)
	assert.Nil(err)
	assert.Equal(1, mdb.executeCalls)
}

func TestQueryExecutor_RunQueryAndScanSingleResult(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockRows{}
	mdb := &mockDb{
		rows: mr,
	}

	qe := NewQueryExecutor(mdb)

	err := qe.RunQueryAndScanSingleResult(context.TODO(), mqb, &mockParser{})
	assert.Nil(err)
	assert.Equal(1, mdb.queryCalls)
	assert.Equal(1, mr.singleValueCalled)
	assert.Equal(1, mr.closeCalled)
}

func TestQueryExecutor_RunQueryAndScanSingleResult_Error(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mdb := &mockDb{
		rows: &mockRows{
			err: errDefault,
		},
	}

	qe := NewQueryExecutor(mdb)

	err := qe.RunQueryAndScanSingleResult(context.TODO(), mqb, &mockParser{})
	assert.Equal(errDefault, err)
}

func TestQueryExecutor_RunQueryAndScanSingleResult_ScanError(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockRows{
		getSingleValueErr: errDefault,
	}
	mdb := &mockDb{
		rows: mr,
	}

	qe := NewQueryExecutor(mdb)

	err := qe.RunQueryAndScanSingleResult(context.TODO(), mqb, &mockParser{})
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbCorruptedData))
	cause := errors.Unwrap(err)
	assert.Equal(errDefault, cause)
	assert.Equal(1, mr.closeCalled)
}

func TestQueryExecutor_RunQueryAndScanAllResults(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockRows{}
	mdb := &mockDb{
		rows: mr,
	}

	qe := NewQueryExecutor(mdb)

	err := qe.RunQueryAndScanAllResults(context.TODO(), mqb, &mockParser{})
	assert.Nil(err)
	assert.Equal(1, mdb.queryCalls)
	assert.Equal(1, mr.allCalled)
	assert.Equal(1, mr.closeCalled)
}

func TestQueryExecutor_RunQueryAndScanAllResults_Error(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mdb := &mockDb{
		rows: &mockRows{
			err: errDefault,
		},
	}

	qe := NewQueryExecutor(mdb)

	err := qe.RunQueryAndScanAllResults(context.TODO(), mqb, &mockParser{})
	assert.Equal(errDefault, err)
}

func TestQueryExecutor_RunQueryAndScanAllResults_ScanError(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockRows{
		getAllErr: errDefault,
	}
	mdb := &mockDb{
		rows: mr,
	}

	qe := NewQueryExecutor(mdb)

	err := qe.RunQueryAndScanAllResults(context.TODO(), mqb, &mockParser{})
	assert.True(errors.IsErrorWithCode(err, errors.ErrDbCorruptedData))
	cause := errors.Unwrap(err)
	assert.Equal(errDefault, cause)
	assert.Equal(1, mr.closeCalled)
}

func TestQueryExecutor_ExecuteQueryAffectingSingleRow(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockResult{
		affectedRows: 1,
	}
	mdb := &mockDb{
		result: mr,
	}

	qe := NewQueryExecutor(mdb)

	err := qe.ExecuteQueryAffectingSingleRow(context.TODO(), mqb)
	assert.Nil(err)
	assert.Equal(1, mdb.executeCalls)
}

func TestQueryExecutor_ExecuteQueryAffectingSingleRow_ExecuteError(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockResult{
		err: errDefault,
	}
	mdb := &mockDb{
		result: mr,
	}

	qe := NewQueryExecutor(mdb)

	err := qe.ExecuteQueryAffectingSingleRow(context.TODO(), mqb)
	assert.Equal(errDefault, err)
}

func TestQueryExecutor_ExecuteQueryAffectingSingleRow_MultipleRowsAffected(t *testing.T) {
	assert := assert.New(t)

	mqb := mockQueryBuilder{}
	mr := &mockResult{
		affectedRows: 2,
	}
	mdb := &mockDb{
		result: mr,
	}

	qe := NewQueryExecutor(mdb)

	err := qe.ExecuteQueryAffectingSingleRow(context.TODO(), mqb)
	assert.True(errors.IsErrorWithCode(err, errors.ErrSqlQueryDidNotAffectSingleRow))
}

type mockQueryBuilder struct {
	buildErr error
}

func (m mockQueryBuilder) Build() (Query, error) {
	return nil, m.buildErr
}

type mockDb struct {
	connectErr    error
	disconnectErr error

	queryCalls int
	queries    []Query
	rows       Rows

	executeCalls int
	executions   []Query
	result       Result
}

func (m *mockDb) Connect(ctx context.Context) error {
	return m.connectErr
}

func (m *mockDb) Disconnect(ctx context.Context) error {
	return m.disconnectErr
}

func (m *mockDb) Query(ctx context.Context, query Query) Rows {
	m.queries = append(m.queries, query)
	m.queryCalls++
	return m.rows
}

func (m *mockDb) Execute(ctx context.Context, query Query) Result {
	m.executions = append(m.executions, query)
	m.executeCalls++
	return m.result
}

type mockRows struct {
	err         error
	closeCalled int
	empty       bool

	singleValueCalled       int
	getSingleValueScannable Scannable
	singleValueScanErr      error
	getSingleValueErr       error

	allCalled       int
	getAllScannable Scannable
	allScanErr      error
	getAllErr       error
}

func (m *mockRows) Err() error {
	return m.err
}

func (m *mockRows) Close() {
	m.closeCalled++
}

func (m *mockRows) Empty() bool {
	return m.empty
}

func (m *mockRows) GetSingleValue(parser RowParser) error {
	m.singleValueCalled++
	if m.getSingleValueScannable != nil {
		m.singleValueScanErr = parser.ScanRow(m.getSingleValueScannable)
	}
	return m.getSingleValueErr
}

func (m *mockRows) GetAll(parser RowParser) error {
	m.allCalled++
	if m.getAllScannable != nil {
		m.allScanErr = parser.ScanRow(m.getAllScannable)
	}
	return m.getAllErr
}

type mockResult struct {
	err                error
	affectedRowsCalled int
	affectedRows       int
}

func (m *mockResult) Err() error {
	return m.err
}

func (m *mockResult) AffectedRows() int {
	m.affectedRowsCalled++
	return m.affectedRows
}
