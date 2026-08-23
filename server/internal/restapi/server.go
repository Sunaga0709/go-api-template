package restapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/environment"
	"github.com/Sunaga0709/go-api-template/internal/logger"
	"github.com/Sunaga0709/go-api-template/internal/restapi/controller"
	"github.com/Sunaga0709/go-api-template/internal/restapi/di"
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
	"github.com/Sunaga0709/go-api-template/internal/restapi/middleware"
)

type Server struct {
	router      *echo.Echo
	port        string
	environment environment.Environment
	shutdownMW  *middleware.ShutdownMiddleware
}

func NewServer(port string, env environment.Environment, lg *logger.Logger, hiddenSuccessLog bool, executorProvider database.ExecutorProvider, txManager database.TxManager) *Server {
	router := echo.New()
	router.HideBanner = true
	router.HidePort = true

	errorHandling := middleware.NewErrorHandlingMiddleware()
	router.HTTPErrorHandler = errorHandling.Handle

	loggerMW := middleware.NewLoggerMiddleware(*lg, hiddenSuccessLog)
	setLocationMW := middleware.NewSetLocationMiddleware()
	panicRecoveryMW := middleware.NewPanicRecoveryMiddleware(*lg)
	shutdownMW := middleware.NewShutdownMiddleware()

	// OpenAPIのStrictミドルウェアではなくEcho全体に適用し、未定義パスや将来追加される非OpenAPIルートも停止中は拒否する。
	router.Use(panicRecoveryMW.Recover)
	router.Use(shutdownMW.RejectDuringShutdown)

	dependency := di.New(executorProvider, txManager)

	controller := controller.NewController(
		dependency.QueryService.CustomerQueryService,
		dependency.Usecase.CustomerUsecase,
		dependency.CommandService.CustomerFavoriteBookCommandService,
		dependency.Usecase.CustomerFavoriteBookUsecase,
		dependency.QueryService.BookQueryService,
	)
	mw := []gen.StrictMiddlewareFunc{ // 先: 内側、後: 外側
		setLocationMW.SetLocation,
		errorHandling.ConvertOpenAPI,
		loggerMW.AccessLogging,
	}

	gen.RegisterHandlers(
		router,
		gen.NewStrictHandler(controller, mw),
	)

	return &Server{
		router:      router,
		port:        port,
		environment: env,
		shutdownMW:  shutdownMW,
	}
}

func (s *Server) Start() error {
	fmt.Printf("start server on [::]:%s (environment: %s)\n", s.port, s.environment.String())
	if err := s.router.Start(listenAddress(s.port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return newError(fmt.Errorf("failed to start server: %w", err))
	}

	return nil
}

func (s *Server) BeginShutdown() {
	s.shutdownMW.BeginShutdown()
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Print("shutting server ... ")
	if err := s.router.Shutdown(ctx); err != nil {
		return newError(fmt.Errorf("failed to shutdown server: %w", err))
	}
	fmt.Println("complete")

	return nil
}

func listenAddress(port string) string {
	port = strings.TrimSpace(port)
	if strings.Contains(port, ":") {
		return port
	}

	return fmt.Sprintf(":%s", port)
}
