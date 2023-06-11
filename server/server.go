package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server interface {
	Run() error
}

type Config struct {
	Port string
}

type server struct {
	app    *fiber.App
	config Config
}

type RouteHandler interface {
	SetupRoutes(fr fiber.Router)
}

func New(config Config, routeHandlers []RouteHandler) Server {
	app := fiber.New()

	app.Use(cors.New(cors.ConfigDefault))

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")

	for _, handler := range routeHandlers {
		handler.SetupRoutes(v1Group)
	}

	s := &server{app: app}

	s.addHealthCheckRoutes()

	return s
}

func (s *server) addHealthCheckRoutes() {
	s.app.Get("/liveness", liveness)
	s.app.Get("/readines", readiness)
}

func liveness(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func readiness(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (s *server) Run() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdownChan
		err := s.app.Shutdown()
		if err != nil {
			log.Println("Error on shutdown gracefully")
		}
	}()
	fmt.Printf("Configuration Port is :%s", s.config.Port)
	return s.app.Listen(fmt.Sprintf(":%s", s.config.Port))
}
