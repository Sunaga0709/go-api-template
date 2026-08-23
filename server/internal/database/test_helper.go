package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

const testDriverName = "database_test_driver"

var (
	testDriverRegisterOnce sync.Once
	testDriverStatesMu     sync.Mutex
	testDriverStates       = map[string]*testDriverState{}
)

type testDriverState struct {
	pingErr       error
	closeCount    atomic.Int32
	beginCount    atomic.Int32
	commitCount   atomic.Int32
	rollbackCount atomic.Int32
}

type testDriver struct{}

type testConn struct {
	state *testDriverState
}

type testTx struct {
	state *testDriverState
}

func (d testDriver) Open(name string) (driver.Conn, error) {
	testDriverStatesMu.Lock()
	state := testDriverStates[name]
	testDriverStatesMu.Unlock()
	if state == nil {
		return nil, fmt.Errorf("unknown test dsn: %s", name)
	}
	return &testConn{state: state}, nil
}

func (c *testConn) Prepare(_ string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not implemented")
}

func (c *testConn) Close() error {
	c.state.closeCount.Add(1)
	return nil
}

func (c *testConn) Begin() (driver.Tx, error) {
	c.state.beginCount.Add(1)
	return &testTx{state: c.state}, nil
}

func (c *testConn) BeginTx(_ context.Context, _ driver.TxOptions) (driver.Tx, error) {
	c.state.beginCount.Add(1)
	return &testTx{state: c.state}, nil
}

func (c *testConn) Ping(_ context.Context) error {
	return c.state.pingErr
}

func (tx *testTx) Commit() error {
	tx.state.commitCount.Add(1)
	return nil
}

func (tx *testTx) Rollback() error {
	tx.state.rollbackCount.Add(1)
	return nil
}

func registerTestDriver() {
	testDriverRegisterOnce.Do(func() {
		sql.Register(testDriverName, testDriver{})
	})
}

func newTestDriverState(t *testing.T, state *testDriverState) string {
	t.Helper()
	registerTestDriver()
	dsn := fmt.Sprintf("%s/%s", t.Name(), t.Name())

	testDriverStatesMu.Lock()
	testDriverStates[dsn] = state
	testDriverStatesMu.Unlock()

	t.Cleanup(func() {
		testDriverStatesMu.Lock()
		delete(testDriverStates, dsn)
		testDriverStatesMu.Unlock()
	})

	return dsn
}

func newTestBunDB(t *testing.T, state *testDriverState) *bun.DB {
	t.Helper()
	dsn := newTestDriverState(t, state)
	sqlDB, err := sql.Open(testDriverName, dsn)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("sqlDB.Close() error = %v", err)
		}
	})

	return bun.NewDB(sqlDB, mysqldialect.New())
}
