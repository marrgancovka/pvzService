package main

import (
	"context"
	"github.com/marrgancovka/pvzService/internal/config"
	"github.com/marrgancovka/pvzService/internal/pkg/db"
	"github.com/marrgancovka/pvzService/internal/pkg/jwter"
	"github.com/marrgancovka/pvzService/internal/pkg/middleware"
	"github.com/marrgancovka/pvzService/internal/pkg/migrations"
	"github.com/marrgancovka/pvzService/internal/pkg/server"
	"github.com/marrgancovka/pvzService/internal/services/auth"
	authHandler "github.com/marrgancovka/pvzService/internal/services/auth/delivery/http"
	authRepository "github.com/marrgancovka/pvzService/internal/services/auth/repo"
	authUsecase "github.com/marrgancovka/pvzService/internal/services/auth/usecase"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	pvzHandler "github.com/marrgancovka/pvzService/internal/services/pvz/delivery/http"
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
			server.NewRouter,
			config.MustLoad,
			fx.Annotate(jwter.New, fx.As(new(auth.JWTer))),

			middleware.NewAuthMiddleware,

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

// TODO: получение пвз
// TODO: тесты 75%
// TODO: интеграционный тест
// TODO: gRPC метод получения пвз
// TODO: добавить прометеус

// TODO: проверить логирование
// TODO: добавить dockerfile + prod.docker-compose + логирование в файл
// TODO: написать makefile
// TODO: написать readme
// TODO: проверить ошибки
// TODO: добавить нужные константы
// TODO: линтер
