package database

import "context"

type TxExecutor interface {
	Executor
	isTxExecutor()
}

type TxManager interface {
	Run(ctx context.Context, fn func(ctx context.Context, executor TxExecutor) error) error
	WithTxExecutor(executor TxExecutor) TxManager
}
