package database

// ExecutorProvider トランザクションを伴わない通常接続のExecutorを提供する
type ExecutorProvider interface {
	Default() Executor
}
