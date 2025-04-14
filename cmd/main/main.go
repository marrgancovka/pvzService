package main

import (
	"context"
	"github.com/marrgancovka/pvzService/internal/config"
	"github.com/marrgancovka/pvzService/internal/pkg/db"
	"github.com/marrgancovka/pvzService/internal/pkg/grpcconn"
	"github.com/marrgancovka/pvzService/internal/pkg/jwter"
	"github.com/marrgancovka/pvzService/internal/pkg/logger"
	"github.com/marrgancovka/pvzService/internal/pkg/metrics"
	"github.com/marrgancovka/pvzService/internal/pkg/middleware"
	"github.com/marrgancovka/pvzService/internal/pkg/servers/mainServer"
	"github.com/marrgancovka/pvzService/internal/pkg/servers/metricsServer"
	"github.com/marrgancovka/pvzService/internal/services/auth"
	authHandler "github.com/marrgancovka/pvzService/internal/services/auth/delivery/http"
	authRepository "github.com/marrgancovka/pvzService/internal/services/auth/repo"
	authUsecase "github.com/marrgancovka/pvzService/internal/services/auth/usecase"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	pvzHandler "github.com/marrgancovka/pvzService/internal/services/pvz/delivery/http"
	pvzRepository "github.com/marrgancovka/pvzService/internal/services/pvz/repo"
	pvzUsecase "github.com/marrgancovka/pvzService/internal/services/pvz/usecase"
	"github.com/marrgancovka/pvzService/migrations"
	"github.com/marrgancovka/pvzService/pkg/builder"
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
			mainServer.NewRouter,
			func() config.ConfigPath {
				return "config/main/config.yaml"
			},
			config.MustLoad,
			fx.Annotate(jwter.New, fx.As(new(auth.JWTer))),

			grpcconn.Provide,

			fx.Annotate(metrics.New, fx.As(new(metrics.Metrics))),
			middleware.NewAuthMiddleware,
			middleware.NewMetricsMiddleware,

			db.NewPostgresPool,
			db.NewPostgresConnect,

			authHandler.NewHandler,
			fx.Annotate(authUsecase.NewUsecase, fx.As(new(auth.Usecase))),
			fx.Annotate(authRepository.NewRepository, fx.As(new(auth.Repository))),

			pvzHandler.NewHandler,
			fx.Annotate(pvzUsecase.NewUsecase, fx.As(new(pvz.Usecase))),
			fx.Annotate(pvzRepository.NewRepository, fx.As(new(pvz.Repository))),
		),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),

		fx.Invoke(
			mainServer.RunServer,
			migrations.RunMigrations,
			metricsServer.RunServer,
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

// TODO: получение пвз

// TODO: написать makefile
// TODO: написать readme
