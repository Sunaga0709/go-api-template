//go:build !integration

package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		env     string
		want    config
		wantErr bool
	}{
		{
			name: "positional arguments",
			args: []string{"./seed", "postgres://user:password@localhost:5432/app"},
			want: config{
				sqlDir:      "./seed",
				databaseURL: "postgres://user:password@localhost:5432/app",
			},
		},
		{
			name: "flags",
			args: []string{"--path", "./seed", "--database-url", "postgres://user:password@localhost:5432/app"},
			want: config{
				sqlDir:      "./seed",
				databaseURL: "postgres://user:password@localhost:5432/app",
			},
		},
		{
			name: "database url from env",
			args: []string{"./seed"},
			env:  "postgres://env-user:password@localhost:5432/app",
			want: config{
				sqlDir:      "./seed",
				databaseURL: "postgres://env-user:password@localhost:5432/app",
			},
		},
		{
			name: "legacy flags",
			args: []string{"-path=./seed", "-db-url=postgres://user:password@localhost:5432/app"},
			want: config{
				sqlDir:      "./seed",
				databaseURL: "postgres://user:password@localhost:5432/app",
			},
		},
		{
			name:    "missing sql dir",
			args:    []string{"--database-url", "postgres://user:password@localhost:5432/app"},
			wantErr: true,
		},
		{
			name:    "missing database url",
			args:    []string{"./seed"},
			wantErr: true,
		},
		{
			name:    "too many positional arguments",
			args:    []string{"./seed", "postgres://user:password@localhost:5432/app", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(databaseURLEnv, tt.env)

			got, err := parseConfig(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeLegacyFlags(t *testing.T) {
	args := []string{
		"-path",
		"./seed",
		"-path=./other-seed",
		"-database-url",
		"postgres://user:password@localhost:5432/app",
		"-database-url=postgres://user:password@localhost:5432/other",
		"-db-url",
		"postgres://user:password@localhost:5432/db",
		"-db-url=postgres://user:password@localhost:5432/legacy",
		"--unchanged",
	}
	want := []string{
		"--path",
		"./seed",
		"--path=./other-seed",
		"--database-url",
		"postgres://user:password@localhost:5432/app",
		"--database-url=postgres://user:password@localhost:5432/other",
		"--db-url",
		"postgres://user:password@localhost:5432/db",
		"--db-url=postgres://user:password@localhost:5432/legacy",
		"--unchanged",
	}

	got := normalizeLegacyFlags(args)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("normalizeLegacyFlags() = %v, want %v", got, want)
	}
}

func Test_sqlFiles(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		dirs      []string
		wantNames []string
		wantErr   bool
	}{
		{
			name: "sorts sql files by numeric prefix",
			files: []string{
				"010_items.sql",
				"001_users.sql",
				"002_orders.SQL",
				"README.md",
			},
			dirs: []string{"003_directory.sql"},
			wantNames: []string{
				"001_users.sql",
				"002_orders.SQL",
				"010_items.sql",
			},
		},
		{
			name: "sorts same prefix by path",
			files: []string{
				"001_b.sql",
				"001_a.sql",
			},
			wantNames: []string{
				"001_a.sql",
				"001_b.sql",
			},
		},
		{
			name:    "no sql files",
			files:   []string{"README.md"},
			wantErr: true,
		},
		{
			name:    "invalid numeric prefix",
			files:   []string{"abc_seed.sql"},
			wantErr: true,
		},
		{
			name:    "missing underscore",
			files:   []string{"001.sql"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, name := range tt.dirs {
				if err := mkdirAll(filepath.Join(dir, name)); err != nil {
					t.Fatal(err)
				}
			}
			for _, name := range tt.files {
				if err := writeFile(filepath.Join(dir, name)); err != nil {
					t.Fatal(err)
				}
			}

			got, err := sqlFiles(dir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("sqlFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			want := make([]string, 0, len(tt.wantNames))
			for _, name := range tt.wantNames {
				want = append(want, filepath.Join(dir, name))
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("sqlFiles() = %v, want %v", got, want)
			}
		})
	}
}

func Test_sqlFiles_readDirError(t *testing.T) {
	_, err := sqlFiles(filepath.Join(t.TempDir(), "not-found"))
	if err == nil {
		t.Fatal("sqlFiles() error = nil, want error")
	}
}

func Test_sqlFileOrder(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    int
		wantErr bool
	}{
		{
			name: "valid numeric prefix",
			arg:  "001_seed.sql",
			want: 1,
		},
		{
			name: "valid zero prefix",
			arg:  "0_seed.sql",
			want: 0,
		},
		{
			name:    "missing underscore",
			arg:     "001.sql",
			wantErr: true,
		},
		{
			name:    "empty prefix",
			arg:     "_seed.sql",
			wantErr: true,
		},
		{
			name:    "non numeric prefix",
			arg:     "abc_seed.sql",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sqlFileOrder(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("sqlFileOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("sqlFileOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mkdirAll(path string) error {
	return os.MkdirAll(path, 0o755)
}

func writeFile(path string) error {
	return os.WriteFile(path, []byte("select 1;"), 0o644)
}
