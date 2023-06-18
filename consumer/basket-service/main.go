package main

import (
	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/basket"
	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/persistence"
	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/product"
	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/config"
	"github.com/eneskzlcn/pact-cdc/httpclient"
	"log"

	"github.com/eneskzlcn/pact-cdc/postgres"
	"github.com/eneskzlcn/pact-cdc/server"
	"github.com/sirupsen/logrus"
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

	repository := persistence.NewPostgresRepository(&persistence.NewPostgresRepositoryOpts{
		DB: db,
		L:  logger,
	})

	httpClient := httpclient.New()

	productClient := product.NewClient(&product.NewClientOpts{
		HTTPClient: httpClient,
		BaseURL:    c.ExternalURL().ProductAPI,
	})

	basketService := basket.NewService(&basket.NewServiceOpts{
		R: repository, L: logger, PC: productClient,
	})

	basketHandler := basket.NewHandler(&basket.NewHandlerOpts{
		S: basketService, L: logger,
	})

	app := server.New(&server.NewServerOpts{
		Port: c.Server().Port,
	}, []server.RouteHandler{
		basketHandler,
	})

	if err := app.Run(); err != nil {
		log.Fatalf("server is closed: %v", err)
	}

}
