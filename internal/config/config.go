package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"log"
	"os"
	"pvzService/internal/pkg/db"
	"pvzService/internal/pkg/jwter"
	"pvzService/internal/pkg/server"
)

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config/config.yaml"`

	HTTPServer server.Config `yaml:"httpServer"`
	DB         db.Config     `yaml:"db"`
	Jwt        jwter.Config  `yaml:"jwt"`
}

type Out struct {
	fx.Out

	HTTPServer server.Config
	DB         db.Config
	Jwt        jwter.Config
}

func MustLoad() Out {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("cannot read .env file: %s", err)
		os.Exit(1)
	}

	if _, err := os.Stat(cfg.ConfigPath); os.IsNotExist(err) {
		log.Printf("config file does not exist: %s", cfg.ConfigPath)
		os.Exit(1)
	}

	if err := cleanenv.ReadConfig(cfg.ConfigPath, &cfg); err != nil {
		log.Printf("cannot read %s: %v", cfg.ConfigPath, err)
		os.Exit(1)
	}

	return Out{
		HTTPServer: cfg.HTTPServer,
		DB:         cfg.DB,
		Jwt:        cfg.Jwt,
	}
}
