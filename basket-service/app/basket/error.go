package basket

import (
	"github.com/eneskzlcn/pact-cdc/cerr"
)

// basket specific errorr

const (
	BasketNotFoundErrCode           cerr.Code = 10100
	ProductNotHasEnoughStockErrCode cerr.Code = 10101
)
