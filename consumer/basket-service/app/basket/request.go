package basket

import "errors"

type CreateBasketRequest struct {
	UserID string `json:"user_id"`
}

func (cbr CreateBasketRequest) Validate() error {
	if cbr.UserID == "" {
		return errors.New("user id can not be empty")
	}

	return nil
}
