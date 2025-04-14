package metricsServer

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

const (
	metricAddr = "localhost:9000"
)

type In struct {
	fx.In

	Logger *slog.Logger
}

func RunServer(in In) {
	metrics := mux.NewRouter().PathPrefix("/metrics").Subrouter()
	metrics.Handle("", promhttp.Handler())

	srv := &http.Server{
		Addr:    metricAddr,
		Handler: metrics,
	}
	go func() {
		in.Logger.Info("starting metricsServer", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}
