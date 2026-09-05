package main

import (
	"os"
	"testing"
)

// setEnvForTest は config / secret に関係する環境変数を一旦すべてクリアした上で、
// env で指定された値のみ設定する。テスト終了時には元の状態へ復元する。
func setEnvForTest(t *testing.T, env map[string]string) {
	t.Helper()
	for _, k := range []string{"SERVER_PORT", "ENVIRONMENT", "LOG_LEVEL", "HIDDEN_SUCCESS_LOG", "DATABASE_URL"} {
		if orig, ok := os.LookupEnv(k); ok {
			t.Cleanup(func() { os.Setenv(k, orig) })
		} else {
			t.Cleanup(func() { os.Unsetenv(k) })
		}
		os.Unsetenv(k)
	}
	for k, v := range env {
		os.Setenv(k, v)
	}
}
