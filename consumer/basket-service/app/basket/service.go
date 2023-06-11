package basket

import "context"

type Service interface {
	CreateBasket(ctx context.Context, req CreateBasketRequest) (*Basket, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateBasket(ctx context.Context, req CreateBasketRequest) (*Basket, error) {
	return nil, nil
}
