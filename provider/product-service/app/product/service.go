package product

import (
	"context"
	"database/sql"
	"errors"
	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/cerr"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service interface {
	GetProductByID(ctx context.Context, id string) (*GetProductResponse, error)
	GetProductsByIDs(
		ctx context.Context, req GetProductsByIDsRequest) (*GetProductsResponse, error)
	CreateProduct(ctx context.Context, req CreateProductRequest) (*CreateProductResponse, error)
}

type service struct {
	logger     *logrus.Logger
	repository Repository
}

type NewServiceOpts struct {
	L *logrus.Logger
	R Repository
}

func NewService(opts *NewServiceOpts) Service {
	return &service{
		logger:     opts.L,
		repository: opts.R,
	}
}

func (s *service) GetProductByID(ctx context.Context, id string) (*GetProductResponse, error) {
	product, err := s.repository.GetProductByID(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.WithField("product_id", id).Errorf("could not get product: %v", err)
		return nil, cerr.Processing()
	}

	if product == nil {
		return nil, cerr.Bag{Code: ProductNotFoundErrCode, Message: "Product not found."}
	}

	return NewGetProductResponse(product), nil
}

func (s *service) GetProductsByIDs(
	ctx context.Context, req GetProductsByIDsRequest) (*GetProductsResponse, error) {
	products := make([]Product, 0, len(req.IDs))
	for _, id := range req.IDs {
		product, err := s.repository.GetProductByID(ctx, id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			s.logger.WithField("product_id", id).Errorf("could not get product: %v", err)
			return nil, cerr.Processing()
		}

		if product == nil {
			return nil, cerr.Bag{Code: OneOrMoreProductsNotFoundErrCode,
				Message: "At least one of given product ids does not exist."}
		}

		products = append(products, *product)
	}

	return NewGetProductsResponse(products), nil
}

func (s *service) CreateProduct(
	ctx context.Context, req CreateProductRequest) (*CreateProductResponse, error) {
	productID := uuid.New().String()
	product, err := s.repository.CreateProduct(ctx, &Product{
		ID:           productID,
		Name:         req.Name,
		Code:         req.Code,
		Color:        req.Color,
		BuyingPrice:  req.BuyingPrice,
		SellingPrice: req.SellingPrice,
		ImageURL:     req.ImageURL,
		Type:         ProductType(req.Type),
		Provider:     req.Provider,
		Creator:      req.Creator,
		Distributor:  req.Distributor,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Errorf("could not create product: %v", err)
		return nil, cerr.Processing()
	}

	return NewCreateProductResponse(product), nil
}
