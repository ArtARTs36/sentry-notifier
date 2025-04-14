package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/handler"
	"github.com/artarts36/sentry-notifier/internal/health"
	"github.com/artarts36/sentry-notifier/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
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

var configCandidates = []string{
	"env://SENTRY_NOTIFIER_CONFIG",
	"sentry-notifier.yaml",
	"sentry-notifier.json",
}

var version = "v1.0.0"

func resolveConfigPath(store storage.Storage) (string, error) {
	for _, candidate := range configCandidates {
		exists, err := store.Exists(candidate)
		if err != nil {
			return "", fmt.Errorf("check exists path %q: %w", candidate, err)
		}

		if exists {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("config not found, scanned in: [%s]", strings.Join(configCandidates, ","))
}

func loadConfig(ctx context.Context) (cfg.Config, error) {
	slog.Info("[main] loading configuration")

	store := storage.Resolve(storage.NewResolver(
		storage.NewFilesystem(),
		map[string]storage.Storage{
			"env://": storage.NewEnv(),
		},
	))

	configPath, err := resolveConfigPath(store)
	if err != nil {
		return cfg.Config{}, fmt.Errorf("resolve path: %w", err)
	}

	slog.Info("[main] config path resolved", slog.String("path", configPath))

	configLoader := newLoader(store)
	config, err := configLoader.Load(ctx, configPath)
	if err != nil {
		return cfg.Config{}, err
	}

	return config, nil
}

func main() {
	setupLogger(slog.LevelDebug)

	slog.Debug("running sentry-notifier", slog.String("version", version))

	ctx, cancel := context.WithCancel(context.Background())

	config, err := loadConfig(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "[main] failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	slog.DebugContext(ctx, "[main] setup log level", slog.String("level", config.Log.Level.String()))
	setupLogger(config.Log.Level.Level())

	slog.
		Info("[main] configuration loaded")

	metricsRegistry := goMetrics.NewDefaultRegistry(goMetrics.Config{
		Namespace: "sentry_notifier",
	})

	metrics.NewAppInfo(metricsRegistry).SetInfo(version, "telegram,mattermost")

	hServer, notifier := app.New(config, metricsRegistry)
	controlServer := registerControl(config, hServer, handler.NewTestHandler(notifier))

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.InfoContext(ctx, fmt.Sprintf("[main] starting main HTTP server on %s", config.HTTP.Addr))

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
		slog.InfoContext(ctx, fmt.Sprintf("[main] starting control HTTP server on %s", config.Control.Addr))

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
		{
			name: "notifier",
			closeFn: func(_ context.Context) error {
				notifier.Close()
				return nil
			},
		},
	}, cancel)

	wg.Wait()
}

type closer struct {
	closeFn func(ctx context.Context) error
	name    string
}

func newLoader(store storage.Storage) *loader.Loader {
	return loader.New(
		store,
		parser.NewResolver(),
		injector.NewComposite([]injector.Injector{
			injector.NewTemplateID(),
			injector.NewNotifyDefaultStrategy(),
		}),
	)
}

func setupLogger(lvl slog.Level) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))

	slog.SetDefault(logger)
}

func registerControl(config cfg.Config, target *app.Server, testHandler *handler.TestHandler) *http.Server {
	const readTimeout = 30 * time.Second

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/health", health.Handler(target.Health))
	mux.Handle("/test", testHandler)

	hServer := &http.Server{
		Addr:        config.Control.Addr,
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
		Info("[main] shutdown..")

	ctx, shCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shCancel()

	for _, cl := range closers {
		slog.Info("[main] closing", slog.String("object", cl.name))

		if err := cl.closeFn(ctx); err != nil {
			slog.
				With(slog.Any("err", err)).
				With(slog.String("object", cl.name)).
				Error("[main] failed to close")
		}
	}

	cancel()
}
