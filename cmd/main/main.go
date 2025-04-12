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
	"pvzService/internal/pkg/jwter"
	"pvzService/internal/pkg/middleware"
	"pvzService/internal/pkg/migrations"
	"pvzService/internal/pkg/server"
	"pvzService/internal/services/auth"
	authHandler "pvzService/internal/services/auth/delivery/http"
	authRepository "pvzService/internal/services/auth/repo"
	authUsecase "pvzService/internal/services/auth/usecase"
	"pvzService/internal/services/pvz"
	pvzHandler "pvzService/internal/services/pvz/delivery/http"
	pvzRepository "pvzService/internal/services/pvz/repo"
	pvzUsecase "pvzService/internal/services/pvz/usecase"
	"pvzService/pkg/builder"
	"pvzService/pkg/logger"
	"syscall"
)

func main() {
	app := fx.New(
		fx.Provide(
			logger.SetupLogger,
			builder.SetupBuilder,
			server.NewRouter,
			config.MustLoad,
			jwter.New,

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
