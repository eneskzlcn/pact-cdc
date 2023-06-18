package product_test

import (
	"fmt"
	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/product"
	"github.com/eneskzlcn/pact-cdc/httpclient"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pact-foundation/pact-go/dsl"
)

var (
	pact          *dsl.Pact
	pactCleanUp   func()
	client        product.Client
	pactServerURL string
)

func TestProductConsumer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Product Consumer Suite")
}

var _ = BeforeSuite(func() {
	pact, pactCleanUp = createPact()

	pactServerURL = fmt.Sprintf("http://localhost:%d", pact.Server.Port)
	client = product.NewClient(&product.NewClientOpts{
		HTTPClient: httpclient.New(),
		BaseURL:    pactServerURL,
	})
})

var _ = AfterSuite(func() {
	defer pactCleanUp()
})

func createPact() (pact *dsl.Pact, cleanUp func()) {
	pact = &dsl.Pact{
		Host:     "localhost",
		Consumer: "BasketService",
		Provider: "ProductService",

		DisableToolValidityCheck: true,
		PactFileWriteMode:        "overwrite",
		LogDir:                   "./pacts/logs",
	}
	//it must be used otherwise it could not create pact file
	pact.Setup(true)

	cleanUp = func() { pact.Teardown() }

	return pact, cleanUp
}
