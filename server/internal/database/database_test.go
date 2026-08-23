package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/uptrace/bun"
)

func TestDatabaseURLConfig(t *testing.T) {
	for _, tt := range []struct {
		name, url, wantDriver, wantDSN string
		wantErr                        bool
	}{{"mysql", "mysql://user:pass@host/db", "mysql", "user:pass@host/db", false}, {"missing separator", "mysql", "", "", true}, {"empty driver", "://dsn", "", "", true}, {"empty dsn", "mysql://", "", "", true}, {"unsupported", "postgres://dsn", "", "", true}} {
		t.Run(tt.name, func(t *testing.T) {
			driverName, dsn, dialect, err := databaseURLConfig(tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("databaseURLConfig() error = %v", err)
			}
			if !tt.wantErr && (driverName != tt.wantDriver || dsn != tt.wantDSN || dialect == nil) {
				t.Errorf("databaseURLConfig() = %q, %q, %v", driverName, dsn, dialect)
			}
		})
	}
}

func TestNewBunDB(t *testing.T) {
	if _, err := newBunDB(context.Background(), "invalid", sql.Open); err == nil {
		t.Error("newBunDB(invalid) error = nil")
	}
	openErr := errors.New("open failed")
	if _, err := newBunDB(context.Background(), "mysql://dsn", func(string, string) (*sql.DB, error) { return nil, openErr }); err == nil {
		t.Error("newBunDB(open error) error = nil")
	}
	state := &testDriverState{}
	dsn := newTestDriverState(t, state)
	db, err := newBunDB(context.Background(), "mysql://test", func(string, string) (*sql.DB, error) { return sql.Open(testDriverName, dsn) })
	if err != nil || db == nil {
		t.Fatalf("newBunDB() = %v, %v", db, err)
	}
	if err := db.Close(); err != nil {
		t.Error(err)
	}
	pingState := &testDriverState{pingErr: errors.New("ping failed")}
	pingDSN := newTestDriverState(t, pingState)
	if _, err := newBunDB(context.Background(), "mysql://test", func(string, string) (*sql.DB, error) { return sql.Open(testDriverName, pingDSN) }); err == nil {
		t.Error("newBunDB(ping error) error = nil")
	}
}

func TestConfigureConnectionPool(t *testing.T) {
	dsn := newTestDriverState(t, &testDriverState{})
	db, err := sql.Open(testDriverName, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	configureConnectionPool(db)
	stats := db.Stats()
	if stats.MaxOpenConnections != dbMaxOpenConns {
		t.Errorf("MaxOpenConnections = %d", stats.MaxOpenConnections)
	}
}

func TestBunExecutorsAndTxManager(t *testing.T) {
	db := newTestBunDB(t, &testDriverState{})
	executor := NewBunExecutor(db)
	if got, err := UnwrapBun(executor); err != nil || got != db {
		t.Errorf("UnwrapBun() = %v, %v", got, err)
	}
	if _, err := UnwrapBun(nil); err == nil {
		t.Error("UnwrapBun(nil) error = nil")
	}
	provider := NewBunExecutorProvider(db)
	if provider.Default() == nil {
		t.Error("Default() = nil")
	}
	manager := NewBunTxManager(db)
	called := false
	if err := manager.Run(context.Background(), func(_ context.Context, tx TxExecutor) error {
		called = true
		if _, err := UnwrapBun(tx); err != nil {
			t.Error(err)
		}
		return nil
	}); err != nil || !called {
		t.Errorf("Run() called=%v err=%v", called, err)
	}
	var captured TxExecutor
	if err := manager.Run(context.Background(), func(_ context.Context, tx TxExecutor) error { captured = tx; return nil }); err != nil {
		t.Fatal(err)
	}
	joined := manager.WithTxExecutor(captured)
	if err := joined.Run(context.Background(), func(_ context.Context, tx TxExecutor) error {
		if tx != captured {
			t.Error("existing transaction not reused")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if _, ok := NewBunTxExecutor(bun.Tx{}).(*BunTxExecutor); !ok {
		t.Error("NewBunTxExecutor() type mismatch")
	}
}

func TestErrorAndTestDriver(t *testing.T) {
	if got := newError(errors.New("cause")).Error(); !strings.Contains(got, "cause") {
		t.Errorf("Error() = %q", got)
	}
	state := &testDriverState{}
	dsn := newTestDriverState(t, state)
	conn, err := (testDriver{}).Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	c := conn.(*testConn)
	if _, prepareErr := c.Prepare("select 1"); prepareErr == nil {
		t.Error("Prepare() error = nil")
	}
	if pingErr := c.Ping(context.Background()); pingErr != nil {
		t.Error(pingErr)
	}
	tx, err := c.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		t.Fatal(commitErr)
	}
	tx, err = c.BeginTx(context.Background(), driver.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		t.Fatal(rollbackErr)
	}
	if closeErr := c.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	if state.beginCount.Load() != 2 || state.commitCount.Load() != 1 || state.rollbackCount.Load() != 1 || state.closeCount.Load() != 1 {
		t.Errorf("driver state = %#v", state)
	}
	if _, err := (testDriver{}).Open("missing"); err == nil {
		t.Error("Open(missing) error = nil")
	}
	registerTestDriver()
	_ = time.Second
}
