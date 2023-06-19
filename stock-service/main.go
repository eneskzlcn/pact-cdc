package main

import (
	"github.com/eneskzlcn/pact-cdc/postgres"
	"github.com/eneskzlcn/pact-cdc/server"
	"github.com/eneskzlcn/pact-cdc/stock-service/app/persistence"
	stock2 "github.com/eneskzlcn/pact-cdc/stock-service/app/stock"
	"github.com/eneskzlcn/pact-cdc/stock-service/config"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	c := config.New()

	db := postgres.New(&postgres.NewPostgresOpts{
		Host:     c.Postgres().Host,
		Port:     c.Postgres().Port,
		DBName:   c.Postgres().DBName,
		Password: c.Postgres().Password,
		Username: c.Postgres().Username,
	})

	logger := logrus.New()

	stockRepository := persistence.NewPostgresRepository(&persistence.NewPostgresRepositoryOpts{
		DB: db,
		L:  logger,
	})

	stockService := stock2.NewService(&stock2.NewServiceOpts{
		R: stockRepository,
		L: logger,
	})

	stockHandler := stock2.NewHandler(&stock2.NewHandlerOpts{
		S: stockService,
		L: logger,
	})

	app := server.New(&server.NewServerOpts{
		Port: c.Server().Port,
	}, []server.RouteHandler{
		stockHandler,
	})

	if err := app.Run(); err != nil {
		log.Fatalf("server is closed: %v", err)
	}
}
