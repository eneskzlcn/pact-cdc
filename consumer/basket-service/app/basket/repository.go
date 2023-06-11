package basket

import "context"

type Repository interface {
	CreateBasket(ctx context.Context, req CreateBasketRequest) (*Basket, error)
}
