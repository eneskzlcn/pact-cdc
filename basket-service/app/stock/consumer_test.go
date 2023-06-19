package stock_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/eneskzlcn/pact-cdc/basket-service/app/stock"
	"github.com/eneskzlcn/pact-cdc/cerr"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pact-foundation/pact-go/dsl"
	"net/http"
)

var _ = Describe("Stock Consumer Test", func() {
	Describe("IsProductAvailableInStockInDesiredQuantity", func() {
		const isProductAvailableInStockPath = "/api/v1/stocks/availability"
		givenProductID, quantity := gofakeit.UUID(), int(gofakeit.Float64Range(0, 10000))

		It("should return error if given no stock information found for given product id", func() {
			pact.AddInteraction().
				Given("i get no stock information found error if no stock information found for given product id").
				UponReceiving("A request for getting available stock information about a product").
				WithRequest(dsl.Request{
					Method:  http.MethodGet,
					Path:    dsl.String(fmt.Sprintf(isProductAvailableInStockPath)),
					Headers: map[string]dsl.Matcher{},
					Body: dsl.StructMatcher{
						"product_id": dsl.Like(givenProductID),
						"quantity":   dsl.Like(quantity),
					},
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusBadRequest,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"code":    30001,
						"message": "No stock information found for given product id.",
					},
				})

			var test = func() error {
				_, err := client.IsProductAvailableInStock(context.Background(), stock.IsProductAvailableInStockRequest{
					ProductID: &givenProductID,
					Quantity:  &quantity,
				})
				return err
			}

			err := pact.Verify(test)
			Expect(err).To(Equal(cerr.Bag{Code: 30001, Message: "No stock information found for given product id."}))
		})

		It("should return false if given product and quantity is not available in stock", func() {
			pact.
				AddInteraction().
				Given("i get product is not available in stock with desired quantity").
				UponReceiving("A request for getting available stock information about a product").
				WithRequest(dsl.Request{
					Method: http.MethodGet,
					Path:   dsl.String(fmt.Sprintf(isProductAvailableInStockPath)),
					Body: dsl.StructMatcher{
						"product_id": dsl.Like(givenProductID),
						"quantity":   dsl.Like(quantity),
					},
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusOK,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"is_available": false,
					},
				})

			var test = func() error {
				_, err := client.IsProductAvailableInStock(context.Background(), stock.IsProductAvailableInStockRequest{
					ProductID: &givenProductID,
					Quantity:  &quantity,
				})
				return err
			}

			err := pact.Verify(test)
			Expect(err).To(BeNil())
		})

		It("should return false if given product and quantity is not available in stock", func() {
			pact.
				AddInteraction().
				Given("i get product is available in stock with desired quantity").
				UponReceiving("A request for getting available stock information about a product").
				WithRequest(dsl.Request{
					Method: http.MethodGet,
					Path:   dsl.String(fmt.Sprintf(isProductAvailableInStockPath)),
					Body: dsl.StructMatcher{
						"product_id": dsl.Like(givenProductID),
						"quantity":   dsl.Like(quantity),
					},
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusOK,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"is_available": true,
					},
				})

			var test = func() error {
				_, err := client.IsProductAvailableInStock(context.Background(), stock.IsProductAvailableInStockRequest{
					ProductID: &givenProductID,
					Quantity:  &quantity,
				})
				return err
			}

			err := pact.Verify(test)
			Expect(err).To(BeNil())
		})
	})
})
