package basket

type Basket struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type BasketProduct struct {
	BasketID  string
	ProductID string
	Quantity  int
}

type BasketProducts []BasketProduct
