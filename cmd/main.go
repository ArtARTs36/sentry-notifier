package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/health"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"

	goMetrics "github.com/artarts36/go-metrics"

	"github.com/artarts36/sentry-notifier/internal/app"
	"github.com/artarts36/sentry-notifier/internal/config/cfg"
	"github.com/artarts36/sentry-notifier/internal/config/injector"
	"github.com/artarts36/sentry-notifier/internal/config/parser"
	"github.com/artarts36/sentry-notifier/internal/config/storage"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artarts36/sentry-notifier/internal/config/loader"
)

func main() {
	setupLogger("debug")

	ctx, cancel := context.WithCancel(context.Background())

	configPath := flag.String("config", "sentry-notifier.yaml", "path to config file")
	if configPath == nil || *configPath == "" {
		slog.ErrorContext(ctx, "[main] missing config path")
		os.Exit(1)
		return
	}

	slog.Info("[main] loading configuration")

	configLoader := newLoader()
	config, err := configLoader.Load(ctx, *configPath)
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			ErrorContext(ctx, "[main] failed to load configuration")

		os.Exit(1)
		return
	}

	setupLogger(config.Log.Level)

	slog.
		Info("[main] configuration loaded")

	hServer := app.New(config, goMetrics.NewDefaultRegistry(goMetrics.Config{
		Namespace: "sentry_notifier",
	}))
	controlServer := registerControl(config, hServer)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.InfoContext(ctx, fmt.Sprintf("[main] starting HTTP server on %s", config.HTTP.Addr))

		err = hServer.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.
				With(slog.Any("err", err)).
				ErrorContext(ctx, "[main] http server listen error")
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.InfoContext(ctx, fmt.Sprintf("[main] starting HTTP server on %s", config.HTTP.Addr))

		err = controlServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.
				With(slog.Any("err", err)).
				ErrorContext(ctx, "[main] metrics server listen error")
			return
		}
	}()

	shutdown([]closer{
		{
			name:    "main-server",
			closeFn: hServer.Shutdown,
		},
		{
			name:    "control-server",
			closeFn: controlServer.Shutdown,
		},
	}, cancel)

	wg.Wait()
}

type closer struct {
	closeFn func(ctx context.Context) error
	name    string
}

func newLoader() *loader.Loader {
	return loader.New(
		storage.NewFilesystem(),
		map[string]parser.Parser{
			"yaml": parser.NewYAML(),
			"yml":  parser.NewYAML(),
		},
		injector.NewComposite([]injector.Injector{
			injector.NewEnv(),
			injector.NewTemplateID(),
			injector.NewNotifyDefaultStrategy(),
		}),
	)
}

func setupLogger(lvl string) {
	level := slog.LevelDebug

	switch lvl {
	case "info":
		level = slog.LevelInfo
	case "warning":
	case "warn":
		level = slog.LevelWarn
	case "error":
	case "err":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}

func registerControl(config cfg.Config, target *app.Server) *http.Server {
	const readTimeout = 30 * time.Second

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/health", health.Handler(target.Health))

	hServer := &http.Server{
		Addr:        config.Metrics.Addr,
		Handler:     mux,
		ReadTimeout: readTimeout,
	}

	return hServer
}

func shutdown(closers []closer, cancel context.CancelFunc) {
	const shutdownTimeout = 30 * time.Second

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	slog.
		With(slog.String("signal", sig.String())).
		Info("shutdown..")

	ctx, shCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shCancel()

	for _, cl := range closers {
		if err := cl.closeFn(ctx); err != nil {
			slog.
				With(slog.Any("err", err)).
				With(slog.String("object", cl.name)).
				Error("failed to close ")
		}
	}

	cancel()
}
