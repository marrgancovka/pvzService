package grpcServer

import "time"

type Config struct {
	Address     string        `yaml:"address" env-default:"localhost:3000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`
}
