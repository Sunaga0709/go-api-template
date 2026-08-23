package database

import (
	"fmt"
	"strings"

	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/schema"
)

func databaseURLConfig(url string) (driverName string, dsn string, dialect schema.Dialect, err error) {
	driverName, dsn, ok := strings.Cut(url, "://")
	if !ok || driverName == "" || dsn == "" {
		return "", "", nil, fmt.Errorf("invalid database url")
	}

	switch driverName {
	case "mysql":
		return driverName, dsn, mysqldialect.New(), nil
	default:
		return "", "", nil, fmt.Errorf("unsupported database driver: %s", driverName)
	}
}
