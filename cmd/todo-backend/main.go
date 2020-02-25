package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/ormanli/todo-backend/internal/app"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	config, err := app.CreateConfig()
	if err != nil {
		logger.Log("config", err)
	}

	r, err := app.NewPostgresRepository(config)
	if err != nil {
		logger.Log("repository", err)
		os.Exit(1)
	}

	s := app.NewService(r, config)
	s = app.LoggingMiddleware(logger)(s)

	h := app.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", config.ServerPort)
		errs <- http.ListenAndServe(config.ServerPort, h)
	}()

	logger.Log("exit", <-errs)
}
