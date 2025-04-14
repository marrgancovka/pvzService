package grpcServer

import (
	pvz "github.com/marrgancovka/pvzService/internal/services/pvz/delivery/grpc"
	"github.com/marrgancovka/pvzService/internal/services/pvz/delivery/grpc/gen"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type In struct {
	fx.In

	Config      Config
	GRPCHandler *pvz.Handler
	Logger      *slog.Logger
}

func RunServer(in In) {
	srv := grpc.NewServer()
	gen.RegisterPVZServiceServer(srv, in.GRPCHandler)
	go func() {
		listener, err := net.Listen("tcp", in.Config.Address)
		if err != nil {
			in.Logger.Error("listen returned err: " + err.Error())
		}
		in.Logger.Info("grpc mainServer started", slog.String("addr", listener.Addr().String()))
		if err = srv.Serve(listener); err != nil {
			in.Logger.Error("serve returned err: " + err.Error())
		}

	}()
}
