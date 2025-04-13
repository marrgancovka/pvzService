package grpcconn

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type In struct {
	fx.In

	Config Config
}

func Provide(in In) (*grpc.ClientConn, error) {
	return grpc.NewClient(in.Config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
