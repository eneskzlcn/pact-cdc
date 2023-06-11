package persistence

import (
	"context"
	"database/sql"

	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/basket"
)

type BasketRepository interface {
	CreateBasket(
		ctx context.Context, req basket.CreateBasketRequest) (*basket.Basket, error)
}

type basketRepository struct {
	db *sql.DB
}

type NewRepositoryOpts struct {
	DB *sql.DB
}

func NewBasketRepository(opts *NewRepositoryOpts) BasketRepository {
	return &basketRepository{
		db: opts.DB,
	}
}

func (r *basketRepository) CreateBasket(
	ctx context.Context, req basket.CreateBasketRequest) (*basket.Basket, error) {
	return nil, nil
}
