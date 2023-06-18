package main

import (
	"github.com/eneskzlcn/pact-cdc/postgres"
	"github.com/eneskzlcn/pact-cdc/provider/product-service/app/persistence"
	"github.com/eneskzlcn/pact-cdc/provider/product-service/app/product"
	"github.com/eneskzlcn/pact-cdc/provider/product-service/config"
	"github.com/eneskzlcn/pact-cdc/server"
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

	var _, _ = db, logger

	productRepository := persistence.NewPostgresRepository(&persistence.NewPostgresRepositoryOpts{
		DB: db,
		L:  logger,
	})

	productService := product.NewService(&product.NewServiceOpts{
		R: productRepository,
		L: logger,
	})

	productHandler := product.NewHandler(&product.NewHandlerOpts{
		S: productService,
		L: logger,
	})

	app := server.New(&server.NewServerOpts{
		Port: c.Server().Port,
	}, []server.RouteHandler{
		productHandler,
	})

	if err := app.Run(); err != nil {
		log.Fatalf("server is closed: %v", err)
	}
}
