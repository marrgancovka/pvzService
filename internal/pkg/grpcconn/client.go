package grpcconn

import (
	"errors"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type In struct {
	fx.In

	Config Config
}

func Provide(in In) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(in.Config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return nil, err
	}
	return conn, err
}
