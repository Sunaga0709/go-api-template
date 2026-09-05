package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/uptrace/bun"

	"github.com/Sunaga0709/go-api-template/internal/database"
)

const databaseURLEnv = "DATABASE_URL"

type config struct {
	sqlDir      string
	databaseURL string
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	conf, err := parseConfig(args)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.NewBunDB(ctx, conf.databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if dbErr := db.Close(); dbErr != nil {
			fmt.Fprintf(os.Stderr, "failed to close database connection: %v\n", dbErr)
		}
	}()

	files, err := sqlFiles(conf.sqlDir)
	if err != nil {
		return err
	}

	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		for _, file := range files {
			fmt.Printf("  --> executing ... %s", file)

			query, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read sql file %q: %w", file, err)
			}
			if strings.TrimSpace(string(query)) == "" {
				continue
			}

			if _, err := tx.ExecContext(ctx, string(query)); err != nil {
				return fmt.Errorf("failed to execute sql file %q: %w", file, err)
			}

			fmt.Println(" ... done")
		}

		return nil
	})
}

func parseConfig(args []string) (config, error) {
	var conf config

	cmd := &cobra.Command{
		Use:           "seeder [sql-dir] [database-url]",
		Short:         "Execute SQL seed files",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args: func(_ *cobra.Command, positional []string) error {
			if len(positional) > 2 {
				return fmt.Errorf("unexpected arguments: %s", strings.Join(positional[2:], " "))
			}
			return nil
		},
		RunE: func(_ *cobra.Command, positional []string) error {
			if conf.sqlDir == "" && len(positional) > 0 {
				conf.sqlDir = positional[0]
			}
			if conf.databaseURL == "" && len(positional) > 1 {
				conf.databaseURL = positional[1]
			}

			if conf.sqlDir == "" {
				return fmt.Errorf("sql directory path is required")
			}
			if conf.databaseURL == "" {
				conf.databaseURL = os.Getenv(databaseURLEnv)
			}
			if conf.databaseURL == "" {
				return fmt.Errorf("database url is required: pass --database-url or set %s", databaseURLEnv)
			}

			return nil
		},
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.SetArgs(normalizeLegacyFlags(args))
	cmd.Flags().StringVar(&conf.sqlDir, "path", "", "directory path containing SQL files")
	cmd.Flags().StringVar(&conf.databaseURL, "database-url", "", "database URL")
	cmd.Flags().StringVar(&conf.databaseURL, "db-url", "", "database URL")

	if err := cmd.Execute(); err != nil {
		return conf, err
	}

	return conf, nil
}

func normalizeLegacyFlags(args []string) []string {
	normalized := make([]string, len(args))
	for i, arg := range args {
		switch {
		case arg == "-path":
			normalized[i] = "--path"
		case strings.HasPrefix(arg, "-path="):
			normalized[i] = "--path=" + strings.TrimPrefix(arg, "-path=")
		case arg == "-database-url":
			normalized[i] = "--database-url"
		case strings.HasPrefix(arg, "-database-url="):
			normalized[i] = "--database-url=" + strings.TrimPrefix(arg, "-database-url=")
		case arg == "-db-url":
			normalized[i] = "--db-url"
		case strings.HasPrefix(arg, "-db-url="):
			normalized[i] = "--db-url=" + strings.TrimPrefix(arg, "-db-url=")
		default:
			normalized[i] = arg
		}
	}
	return normalized
}

func sqlFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read sql directory %q: %w", dir, err)
	}

	type sqlFile struct {
		path  string
		order int
	}

	sqlFiles := make([]sqlFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".sql" {
			continue
		}

		order, err := sqlFileOrder(entry.Name())
		if err != nil {
			return nil, err
		}

		sqlFiles = append(sqlFiles, sqlFile{
			path:  filepath.Join(dir, entry.Name()),
			order: order,
		})
	}

	if len(sqlFiles) == 0 {
		return nil, fmt.Errorf("sql files not found in %q", dir)
	}

	sort.Slice(sqlFiles, func(i, j int) bool {
		if sqlFiles[i].order == sqlFiles[j].order {
			return sqlFiles[i].path < sqlFiles[j].path
		}
		return sqlFiles[i].order < sqlFiles[j].order
	})

	files := make([]string, 0, len(sqlFiles))
	for _, sqlFile := range sqlFiles {
		files = append(files, sqlFile.path)
	}

	return files, nil
}

func sqlFileOrder(name string) (int, error) {
	prefix, _, ok := strings.Cut(name, "_")
	if !ok {
		return 0, fmt.Errorf("sql file %q must start with a numeric prefix followed by _", name)
	}

	order, err := strconv.Atoi(prefix)
	if err != nil {
		return 0, fmt.Errorf("sql file %q has invalid numeric prefix %q: %w", name, prefix, err)
	}

	return order, nil
}
