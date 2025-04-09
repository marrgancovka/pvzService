package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"os"
	"os/signal"
	"pvzService/internal/config"
	"pvzService/internal/pkg/db"
	"pvzService/internal/pkg/migrations"
	"pvzService/internal/pkg/server"
	"pvzService/pkg/logger"
	"syscall"
)

func main() {
	app := fx.New(
		fx.Provide(
			logger.SetupLogger,
			config.MustLoad,

			db.NewPostgresPool,
			db.NewPostgresConnect,
		),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),

		fx.Invoke(
			server.RunServer,
			migrations.RunMigrations,
		),
	)

	ctx := context.Background()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	<-stop
	app.Stop(ctx)
}
