package tests

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/marrgancovka/pvzService/internal/config"
	"github.com/marrgancovka/pvzService/internal/pkg/db"
	"github.com/marrgancovka/pvzService/internal/pkg/jwter"
	"github.com/marrgancovka/pvzService/internal/pkg/servers/mainServer"
	"time"
)

const (
	testURL = "0.0.0.0:8085"
	dbName  = "test_db"
	dbUser  = "test"
	dbPass  = "123"
	dbPort  = 5433
	dbHost  = "localhost"
)

func Config() config.Out {

	return config.Out{
		HTTPServer: mainServer.Config{
			Address:           testURL,
			Timeout:           time.Second * 10,
			IdleTimeout:       time.Second * 10,
			ReadHeaderTimeout: time.Second * 10,
		},
		DB: db.Config{
			DB:             dbName,
			User:           dbUser,
			Password:       dbPass,
			Host:           dbHost,
			Port:           dbPort,
			ConnectTimeout: time.Minute,
		},
		Jwt: jwter.Config{
			ExpirationTime: time.Hour,
			KeyJWT:         []byte("dknsnslk"),
		},
	}
}
