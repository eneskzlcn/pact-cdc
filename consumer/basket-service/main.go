package main

import (
	"log"

	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/basket"
	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/persistence"

	"github.com/eneskzlcn/pact-cdc/postgres"
	"github.com/eneskzlcn/pact-cdc/server"
)

func main() {
	db := postgres.New(postgres.Config{
		Host:     "localhost",
		Port:     "5432",
		DBName:   "pact-cdc",
		Password: "pact-cdc",
		Username: "pact-cdc",
	})

	basketRepo := persistence.NewBasketRepository(&persistence.NewRepositoryOpts{
		DB: db,
	})

	basketService := basket.NewService(basketRepo)

	basketHandler := basket.NewHandler(basketService)

	app := server.New(server.Config{
		Port: "9000",
	}, []server.RouteHandler{
		basketHandler,
	})

	if err := app.Run(); err != nil {
		log.Fatalf("server is closed: %v", err)
	}

}
