package basket

import "github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/product"

type BasketResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Products  []product.Product `json:"products,omitempty"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

func NewBasketResponse(basket *Basket, products []product.Product) *BasketResponse {
	if basket == nil {
		return nil
	}

	return &BasketResponse{
		ID:        basket.ID,
		UserID:    basket.UserID,
		CreatedAt: basket.CreatedAt.Format(layoutISO),
		UpdatedAt: basket.UpdatedAt.Format(layoutISO),
		Products:  products,
	}
}
