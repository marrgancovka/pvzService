package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/marrgancovka/pvzService/internal/pkg/db"
	"github.com/marrgancovka/pvzService/internal/pkg/grpcconn"
	"github.com/marrgancovka/pvzService/internal/pkg/grpcserver"
	"github.com/marrgancovka/pvzService/internal/pkg/jwter"
	"github.com/marrgancovka/pvzService/internal/pkg/server"
	"go.uber.org/fx"
	"log"
	"os"
)

type Config struct {
	//ConfigPath string `env:"CONFIG_PATH" env-default:"config/config.yaml"`

	HTTPServer    server.Config     `yaml:"httpServer"`
	GRPCServer    grpcserver.Config `yaml:"grpcServer"`
	PvzGRPCClient grpcconn.Config   `yaml:"pvzGRPCClient"`
	DB            db.Config         `yaml:"db"`
	Jwt           jwter.Config      `yaml:"jwt"`
}

type ConfigPath string

type In struct {
	fx.In

	Path ConfigPath
}

type Out struct {
	fx.Out

	HTTPServer    server.Config
	GRPCServer    grpcserver.Config
	PvzGRPCClient grpcconn.Config
	DB            db.Config
	Jwt           jwter.Config
}

func MustLoad(in In) Out {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("cannot read .env file: %s", err)
		os.Exit(1)
	}

	if _, err := os.Stat(string(in.Path)); os.IsNotExist(err) {
		log.Printf("config file does not exist: %s", string(in.Path))
		os.Exit(1)
	}

	if err := cleanenv.ReadConfig(string(in.Path), &cfg); err != nil {
		log.Printf("cannot read %s: %v", string(in.Path), err)
		os.Exit(1)
	}

	return Out{
		HTTPServer:    cfg.HTTPServer,
		GRPCServer:    cfg.GRPCServer,
		PvzGRPCClient: cfg.PvzGRPCClient,
		DB:            cfg.DB,
		Jwt:           cfg.Jwt,
	}
}
