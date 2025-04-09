package jwter

import "time"

type Config struct {
	ExpirationTime time.Duration `yaml:"accessExpirationTime" env-default:"24h"`
	KeyJWT         []byte        `env:"JWT_SECRET" env-required:"true"`
}
