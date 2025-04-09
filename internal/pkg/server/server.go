package server

import (
	"go.uber.org/fx"
	"net/http"
)

type Params struct {
	fx.In

	Config Config
}

func RunServer(params Params) {
	srv := &http.Server{
		Addr:              params.Config.Address,
		ReadHeaderTimeout: params.Config.ReadHeaderTimeout,
		IdleTimeout:       params.Config.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}
