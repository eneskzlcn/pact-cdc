package basket

import "github.com/gofiber/fiber/v2"

type Handler interface {
	SetupRoutes(fr fiber.Router)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) CreateBasket(c *fiber.Ctx) error {
	return nil
}

func (h *handler) SetupRoutes(fr fiber.Router) {
	basketGroup := fr.Group("/basket")

	basketGroup.Post("/", h.CreateBasket)
}
