package grpcconn

type Config struct {
	Address string `yaml:"address" env-default:"localhost:3000"`
}
