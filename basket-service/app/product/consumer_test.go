package product_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/eneskzlcn/pact-cdc/basket-service/app/product"
	"github.com/eneskzlcn/pact-cdc/cerr"
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

	/*Describe("GetProductsByIDs", func() {
		const getProductsByIDsPath = "/api/v1/products/bulk"

		givenProductIDs := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

		givenReq := product.GetProductByIDsRequest{
			IDs: givenProductIDs,
		}

		It("should return validation error if no product id is given", func() {
			pact.AddInteraction().
				Given("i get validation error when no product id is given").
				UponReceiving("A request for get products").
				WithRequest(dsl.Request{
					Method: http.MethodGet,
					Path:   dsl.String(getProductsByIDsPath),
					Body:   nil,
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusBadRequest,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"code":    20002,
						"message": "At least one product id must be given.",
					},
				})

			var test = func() error {
				_, err := client.GetProductsByIDs(context.Background(), givenReq)
				return err
			}

			err := pact.Verify(test)
			Expect(err).To(Equal(cerr.Bag{Code: 20002, Message: "At least one product id must be given."}))
		})

		It("should return error if any of product with given ids does not exist", func() {
			pact.AddInteraction().
				Given("i get product not found error when the one of product with given id does not exists").
				UponReceiving("A request for get products contains at least one not exist product id").
				WithRequest(dsl.Request{
					Method: http.MethodGet,
					Path:   dsl.String(getProductsByIDsPath),
					Body: dsl.StructMatcher{
						"ids": dsl.Like(givenProductIDs),
					},
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusBadRequest,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"code":    20003,
						"message": "At least one of given product ids does not exist.",
					},
				})

			var test = func() error {
				_, err := client.GetProductsByIDs(context.Background(), givenReq)
				return err
			}

			err := pact.Verify(test)
			Expect(err).To(Equal(cerr.Bag{Code: 20003, Message: "At least one of given product ids does not exist."}))
		})

		It("should return error if all of products with given ids does not exist", func() {
			pact.AddInteraction().
				Given("i get product not found error when the product with given id does not exists").
				UponReceiving("A request for get products contains at least one not exist product id").
				WithRequest(dsl.Request{
					Method: http.MethodGet,
					Path:   dsl.String(getProductsByIDsPath),
					Body: dsl.StructMatcher{
						"ids": dsl.Like(givenProductIDs),
					},
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusBadRequest,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"code":    20004,
						"message": "None of given product ids exist.",
					},
				})

			var test = func() error {
				_, err := client.GetProductsByIDs(context.Background(), givenReq)
				return err
			}

			err := pact.Verify(test)
			Expect(err).To(Equal(cerr.Bag{Code: 20004, Message: "None of given product ids exist."}))
		})

		It("should return products if given product ids exists", func() {
			givenProducts := []product.Product{
				{
					ID:        givenProductIDs[0],
					Name:      gofakeit.Name(),
					Code:      gofakeit.Word(),
					Color:     gofakeit.Color(),
					CreatedAt: gofakeit.Date(),
					UpdatedAt: gofakeit.Date(),
					Price:     gofakeit.Price(10, 100),
					ImageURL:  gofakeit.ImageURL(200, 100),
					Type:      gofakeit.Word(),
				},
				{
					ID:        givenProductIDs[1],
					Name:      gofakeit.Name(),
					Code:      gofakeit.Word(),
					Color:     gofakeit.Color(),
					CreatedAt: gofakeit.Date(),
					UpdatedAt: gofakeit.Date(),
					Price:     gofakeit.Price(10, 100),
					ImageURL:  gofakeit.ImageURL(200, 100),
					Type:      gofakeit.Word(),
				},
				{
					ID:        givenProductIDs[2],
					Name:      gofakeit.Name(),
					Code:      gofakeit.Word(),
					Color:     gofakeit.Color(),
					CreatedAt: gofakeit.Date(),
					UpdatedAt: gofakeit.Date(),
					Price:     gofakeit.Price(10, 100),
					ImageURL:  gofakeit.ImageURL(200, 100),
					Type:      gofakeit.Word(),
				},
			}

			pact.AddInteraction().
				Given("i get products with given ids").
				UponReceiving("A request for get products with given ids").
				WithRequest(dsl.Request{
					Method: http.MethodGet,
					Path:   dsl.String(getProductsByIDsPath),
					Body: dsl.StructMatcher{
						"ids": dsl.Like(givenProductIDs),
					},
				}).
				WillRespondWith(dsl.Response{
					Status: http.StatusOK,
					Headers: map[string]dsl.Matcher{
						fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
					},
					Body: dsl.StructMatcher{
						"products": dsl.EachLike(dsl.StructMatcher{
							"id":        dsl.Like(givenProducts[0].ID),
							"name":      dsl.Like(givenProducts[0].Name),
							"code":      dsl.Like(givenProducts[0].Code),
							"color":     dsl.Like(givenProducts[0].Color),
							"createdAt": dsl.Like(givenProducts[0].CreatedAt),
							"updatedAt": dsl.Like(givenProducts[0].UpdatedAt),
							"price":     dsl.Like(givenProducts[0].Price),
							"imageURL":  dsl.Like(givenProducts[0].ImageURL),
							"type":      dsl.Like(givenProducts[0].Type),
						}, len(givenProducts)),
					},
				})

			var test = func() error {
				_, err := client.GetProductsByIDs(context.Background(), givenReq)
				return err
			}

			err := pact.Verify(test)
			Expect(err).To(BeNil())
		})
	})*/
})
