package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"

	"loan-ranger/internal/app"
	"loan-ranger/internal/pkg"
	"loan-ranger/internal/pkg/config"
)

type Server struct {
	E    *echo.Echo
	conf *config.Config
	opt  *pkg.Options
}

func NewHTTPServer() Server {
	var (
		conf   = config.GetConfig()
		appCtx = app.ContextApp{Config: conf}
	)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	s := Server{
		E:    echo.New(),
		conf: &conf,
		opt: &pkg.Options{
			Config: conf,
			DB:     appCtx.GetDB(),
			Bucket: appCtx.GetS3BucketClient(),
		},
	}
	s.initRouter()

	return s
}

func (s *Server) Start() {
	go func() {
		if err := s.E.Start(fmt.Sprintf(":%s", "5000")); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.E.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}

func (s *Server) Stop() {
	if err := s.E.Shutdown(context.Background()); err != nil {
		s.E.Logger.Fatal(err)
	}
	if s.opt.DB != nil {
		_ = s.opt.DB.Close()
	}
	slog.Info("service gracefully shutdown")
}
