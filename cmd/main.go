package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/app"
	"github.com/artarts36/sentry-notifier/internal/config/injector"
	"github.com/artarts36/sentry-notifier/internal/config/parser"
	"github.com/artarts36/sentry-notifier/internal/config/storage"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
			With(slog.String("err", err.Error())).
			ErrorContext(ctx, "[main] failed to load configuration")

		os.Exit(1)
		return
	}

	setupLogger(config.Log.Level)

	slog.
		Info("[main] configuration loaded")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	hServer := app.New(config)

	go func() {
		slog.InfoContext(ctx, fmt.Sprintf("[main] starting HTTP server on %s", config.HTTP.Addr))

		err = hServer.Run()
		if err != nil {
			slog.
				With(slog.String("err", err.Error())).
				ErrorContext(ctx, "[main] http server listen error")

			os.Exit(1)
			return
		}
	}()

	<-done

	slog.InfoContext(ctx, "[main] stopping http server")

	err = hServer.Shutdown(ctx)
	if err != nil {
		slog.
			With(slog.String("err", err.Error())).
			ErrorContext(ctx, "[main] failed to shutdown http server")
	}

	slog.Info("[main] canceling root context")

	cancel()
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
