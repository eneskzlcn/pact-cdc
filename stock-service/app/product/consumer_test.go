package product_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/eneskzlcn/pact-cdc/cerr"
	"github.com/eneskzlcn/pact-cdc/stock-service/app/product"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pact-foundation/pact-go/dsl"
	"net/http"
)

var _ = Describe("Product Consumer Test", func() {
	Describe("GetProductByID", func() {
		const getProductByIDPath = "/api/v1/products/%s"
		givenProductID := gofakeit.UUID()

		It("should return product not found error if product with given id does not exist", func() {
			pact.AddInteraction().
				Given("i get product not found error when the product with given id does not exists").
				UponReceiving("A request for product with a non exist product id").
				WithRequest(dsl.Request{
					Method:  http.MethodGet,
					Path:    dsl.String(fmt.Sprintf(getProductByIDPath, givenProductID)),
					Headers: map[string]dsl.Matcher{},
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusBadRequest,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"code":    20001,
						"message": "Product not found.",
					},
				})

			var test = func() error {
				_, err := client.GetProductByID(context.Background(), givenProductID)
				return err
			}

			err := pact.Verify(test)
			Expect(err).To(Equal(cerr.Bag{Code: 20001, Message: "Product not found."}))
		})

		It("should return product if product with given id exists", func() {
			givenProduct := product.Product{
				ID:        givenProductID,
				Name:      gofakeit.Name(),
				Code:      gofakeit.Word(),
				Color:     gofakeit.Color(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
				Price:     gofakeit.Price(10, 100),
				ImageURL:  gofakeit.ImageURL(200, 100),
				Type:      gofakeit.Word(),
			}

			pact.AddInteraction().
				Given("i get product with given id").
				UponReceiving("A request for product with a exist product id").
				WithRequest(dsl.Request{
					Method: http.MethodGet,
					Path:   dsl.String(fmt.Sprintf(getProductByIDPath, givenProductID)),
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusOK,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"id":         givenProduct.ID,
						"name":       dsl.Like(givenProduct.Name),
						"code":       dsl.Like(givenProduct.Code),
						"color":      dsl.Like(givenProduct.Color),
						"created_at": dsl.Like(givenProduct.CreatedAt),
						"updated_at": dsl.Like(givenProduct.UpdatedAt),
						"price":      dsl.Like(givenProduct.Price),
						"image_url":  dsl.Like(givenProduct.ImageURL),
						"type":       dsl.Like(givenProduct.Type),
					},
				})

			var test = func() error {
				_, err := client.GetProductByID(context.Background(), givenProductID)
				return err
			}

			err := pact.Verify(test)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
