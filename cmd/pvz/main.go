package main

import (
	"context"
	"github.com/marrgancovka/pvzService/internal/config"
	"github.com/marrgancovka/pvzService/internal/pkg/db"
	"github.com/marrgancovka/pvzService/internal/pkg/grpcserver"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	"github.com/marrgancovka/pvzService/internal/services/pvz/delivery/grpc"
	pvzRepository "github.com/marrgancovka/pvzService/internal/services/pvz/repo"
	pvzUsecase "github.com/marrgancovka/pvzService/internal/services/pvz/usecase"
	"github.com/marrgancovka/pvzService/pkg/builder"
	"github.com/marrgancovka/pvzService/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fx.New(
		fx.Provide(
			logger.SetupLogger,
			builder.SetupBuilder,
			func() config.ConfigPath {
				return "config/pvz/config.yaml"
			},
			config.MustLoad,

			db.NewPostgresPool,
			grpc.NewHandler,
			fx.Annotate(pvzUsecase.NewUsecase, fx.As(new(pvz.Usecase))),
			fx.Annotate(pvzRepository.NewRepository, fx.As(new(pvz.Repository))),
		),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),

		fx.Invoke(
			grpcserver.RunServer,
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
