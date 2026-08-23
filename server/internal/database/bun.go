package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQLドライバ
	"github.com/uptrace/bun"
)

func NewBunDB(ctx context.Context, url string) (*bun.DB, error) {
	return newBunDB(ctx, url, sql.Open)
}

func newBunDB(
	ctx context.Context,
	url string,
	open func(driverName string, dataSourceName string) (*sql.DB, error),
) (*bun.DB, error) {
	driverName, dsn, dialect, err := databaseURLConfig(url)
	if err != nil {
		return nil, newError(err)
	}

	sqlDB, err := open(driverName, dsn)
	if err != nil {
		return nil, newError(fmt.Errorf("failed to open database: %w", err))
	}
	configureConnectionPool(sqlDB)

	ctx, cancel := context.WithTimeout(ctx, dbPingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, newError(fmt.Errorf("failed to connect database: %w", err))
	}

	return bun.NewDB(sqlDB, dialect), nil
}

// BunExecutor 通常接続のBun実行主体を保持する
type BunExecutor struct {
	db *bun.DB
}

func NewBunExecutor(db *bun.DB) Executor {
	return &BunExecutor{db: db}
}

func (b *BunExecutor) isExecutor() {}

// BunTxExecutor トランザクション接続のBun実行主体を保持する
type BunTxExecutor struct {
	tx bun.Tx
}

func NewBunTxExecutor(tx bun.Tx) TxExecutor {
	return &BunTxExecutor{tx: tx}
}

func (b *BunTxExecutor) isExecutor() {}

func (b *BunTxExecutor) isTxExecutor() {}

// UnwrapBun ExecutorをBun用Executorに型アサーションし、bun.IDBを取り出す
func UnwrapBun(executor Executor) (bun.IDB, error) {
	switch e := executor.(type) {
	case *BunExecutor:
		return e.db, nil
	case *BunTxExecutor:
		return e.tx, nil
	default:
		return nil, newError(fmt.Errorf("invalid bun executor: %T", executor))
	}
}

// BunExecutorProvider 通常接続のExecutorを提供する
type BunExecutorProvider struct {
	db *bun.DB
}

var _ ExecutorProvider = (*BunExecutorProvider)(nil)

func NewBunExecutorProvider(db *bun.DB) *BunExecutorProvider {
	return &BunExecutorProvider{db: db}
}

func (b *BunExecutorProvider) Default() Executor {
	return NewBunExecutor(b.db)
}

// BunTxManager トランザクション境界を管理する
type BunTxManager struct {
	db         *bun.DB
	txExecutor TxExecutor
}

func NewBunTxManager(db *bun.DB) TxManager {
	return &BunTxManager{db: db}
}

// WithTxExecutor 既存Txに参加するTxManagerを生成する
func (b *BunTxManager) WithTxExecutor(txExecutor TxExecutor) TxManager {
	return &BunTxManager{
		db:         b.db,
		txExecutor: txExecutor,
	}
}

// Run トランザクションを開始し、TxExecutorを渡してfnを実行する。
// 既存TxExecutorを保持している場合は、新しいTxを開始せず、そのTxExecutorを使用する。
func (b *BunTxManager) Run(
	ctx context.Context,
	fn func(ctx context.Context, executor TxExecutor) error,
) error {
	if b.txExecutor != nil {
		// NOTE: 呼び出し元でエラー分類を可能にするため、ラップせずそのまま返却する
		return fn(ctx, b.txExecutor)
	}

	// NOTE: 呼び出し元でエラー分類を可能にするため、ラップせずそのまま返却する
	return b.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(ctx, NewBunTxExecutor(tx))
	})
}
