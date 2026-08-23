package database

import (
	"database/sql"
	"time"
)

// 最大5台のインスタンスにスケールアウトする前提で、
// DBへの総接続数が最大100接続に収まるように1インスタンスあたり20接続に制限
const (
	dbMaxOpenConns    = 20
	dbMaxIdleConns    = 20
	dbConnMaxLifetime = 5 * time.Minute
	dbConnMaxIdleTime = 1 * time.Minute
	dbPingTimeout     = 5 * time.Second
)

func configureConnectionPool(db *sql.DB) {
	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)
	db.SetConnMaxLifetime(dbConnMaxLifetime)
	db.SetConnMaxIdleTime(dbConnMaxIdleTime)
}
