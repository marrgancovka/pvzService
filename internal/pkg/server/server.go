package server

import (
	"errors"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

type Params struct {
	fx.In

	Config Config
	Router *Router
	Logger *slog.Logger
}

func RunServer(params Params) {
	srv := &http.Server{
		Addr:              params.Config.Address,
		Handler:           params.Router.handler,
		ReadHeaderTimeout: params.Config.ReadHeaderTimeout,
		IdleTimeout:       params.Config.IdleTimeout,
	}
	go func() {
		params.Logger.Info("starting server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}
