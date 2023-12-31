package basket

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/cerr"
	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/product"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreateBasket(ctx context.Context, req CreateBasketRequest) (*GetBasketResponse, error)
	AddProductToBasket(ctx context.Context, req AddProductToBasketRequest) (*GetBasketResponse, error)
	GetBasketByID(ctx context.Context, basketID string) (*GetBasketResponse, error)
	AddBulkProductToBasket(ctx context.Context, req AddBulkProductToBasketRequest) (*GetBasketResponse, error)
}

type service struct {
	repo          Repository
	logger        *logrus.Logger
	productClient product.Client
}

type NewServiceOpts struct {
	R  Repository
	L  *logrus.Logger
	PC product.Client
}

func NewService(opts *NewServiceOpts) Service {
	return &service{
		repo:          opts.R,
		logger:        opts.L,
		productClient: opts.PC,
	}
}

func (s *service) CreateBasket(
	ctx context.Context, req CreateBasketRequest) (*GetBasketResponse, error) {
	basketID := uuid.New().String()

	basket, err := s.repo.CreateBasket(ctx, &Basket{
		ID:     basketID,
		UserID: req.UserID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("could not create basket: %v", err)
		return nil, cerr.Processing()
	}

	return NewBasketResponse(basket, nil), nil
}

func (s *service) AddProductToBasket(
	ctx context.Context, req AddProductToBasketRequest) (*GetBasketResponse, error) {
	basket, err := s.repo.GetBasketByID(ctx, req.BasketID)
	if err != nil || basket == nil || basket.UserID != req.UserID {
		s.logger.WithField("basket_id", req.BasketID).Error("could not found basket: %v", err)
		return nil, cerr.Bag{Code: BasketNotFoundErrCode, Message: "basket not found"}
	}

	basket, err = s.repo.AddProductToBasket(ctx, &Product{
		ID:       req.ProductID,
		Quantity: req.Quantity,
		BasketID: basket.ID,
	})
	if err != nil {
		s.logger.WithField("basket_id", req.BasketID).Error("could not add product to basket: %v", err)
		return nil, cerr.Processing()
	}

	productIDs := getIDsOfProducts(basket.Products)

	products, err := s.getProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	return NewBasketResponse(basket, products), nil
}

func (s *service) GetBasketByID(
	ctx context.Context, basketID string) (*GetBasketResponse, error) {
	basket, err := s.repo.GetBasketByID(ctx, basketID)
	if err != nil {
		s.logger.WithField("basket_id", basketID).Error("could not found basket: %v", err)
		return nil, cerr.Processing()
	}

	productIDs := getIDsOfProducts(basket.Products)

	products, err := s.getProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}
	
	return NewBasketResponse(basket, products), nil
}

func (s *service) AddBulkProductToBasket(
	ctx context.Context, req AddBulkProductToBasketRequest) (*GetBasketResponse, error) {
	basket, err := s.repo.GetBasketByID(ctx, req.BasketID)
	if err != nil || basket == nil || basket.UserID != req.UserID {
		s.logger.WithField("basket_id", req.BasketID).Error("could not found basket: %v", err)
		return nil, cerr.Bag{Code: BasketNotFoundErrCode, Message: "basket not found"}
	}

	for _, prod := range req.Products {
		_, err := s.repo.AddProductToBasket(ctx, &Product{
			ID:       prod.ID,
			Quantity: prod.Quantity,
			BasketID: basket.ID,
		})
		if err != nil {
			s.logger.WithField("basket_id", req.BasketID).
				WithField("product_id", prod.ID).Error("could not add product to basket: %v", err)
			return nil, cerr.Processing()
		}
	}

	return s.GetBasketByID(ctx, basket.ID)
}

func (s *service) getProductByID(ctx context.Context, productID string) (*product.Product, error) {
	prod, err := s.productClient.GetProductByID(ctx, productID)
	if err != nil {
		s.logger.WithField("product_id", productID).Error("could not get product from product service: %v", err)
		return nil, cerr.Processing()
	}

	return prod, nil
}

func (s *service) getProductsByIDs(ctx context.Context, productIDs []string) ([]product.Product, error) {
	products, err := s.productClient.GetProductsByIDs(ctx, product.GetProductByIDsRequest{
		IDs: productIDs,
	})
	if err != nil {
		s.logger.Error("could not get products from product api: %v", err)
		return nil, cerr.Processing()
	}

	return products, nil
}
