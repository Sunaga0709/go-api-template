package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/environment"
	"github.com/Sunaga0709/go-api-template/internal/logger"
	"github.com/Sunaga0709/go-api-template/internal/restapi"
)

const gracefulShutdownDelay = 5

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var err error

	var (
		conf config
		secr secret
		env  environment.Environment
		lg   *logger.Logger
	)
	{
		{
			fmt.Print("loading config ... ")
			conf, err = loadConfig()
			if err != nil {
				return err
			}
			fmt.Println("complete")
		}

		{
			fmt.Print("loading secret ... ")
			secr, err = loadSecret()
			if err != nil {
				return err
			}
			fmt.Println("complete")
		}

		env, err = environment.NewEnvironment(conf.Environment)
		if err != nil {
			return err
		}
		lg, err = logger.NewLogger(conf.LogLevel)
		if err != nil {
			return err
		}
	}

	var (
		executorProvider database.ExecutorProvider
		txManager        database.TxManager
	)
	{
		fmt.Print("connecting database ... ")
		bunDB, err := database.NewBunDB(ctx, secr.DatabaseURL)
		if err != nil {
			return err
		}
		defer func() {
			if err := bunDB.Close(); err != nil {
				fmt.Printf("failed to close database connection: %v\n", err)
			}
		}()
		fmt.Println("complete")

		executorProvider = database.NewBunExecutorProvider(bunDB)
		txManager = database.NewBunTxManager(bunDB)
	}

	server := restapi.NewServer(conf.Port, env, lg, conf.HiddenSuccessLog, executorProvider, txManager)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	select {
	case <-ctx.Done():
		server.BeginShutdown()
		fmt.Printf("received shutdown signal; stop accepting requests and wait %ds ", gracefulShutdownDelay)
		for range gracefulShutdownDelay {
			time.Sleep(time.Second)
			fmt.Printf(".")
		}
		fmt.Println()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
	case err := <-errCh:
		if err != nil {
			return err
		}
	}

	return nil
}
