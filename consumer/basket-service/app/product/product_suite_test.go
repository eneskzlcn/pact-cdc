package product_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pact-foundation/pact-go/dsl"
)

const (
	productAPIExternalURL = "http://product-api"
	pactBrokerAddress     = "http://localhost:9292"
)

var (
//pact              *dsl.Pact
//pactCleanup       func()
//client            product.Client
//mockConfigManager *config.MockManager
)

func TestProductConsumer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Product Consumer Suite")
}

var _ = BeforeSuite(func() {
	//mockCtrl := gomock.NewController(GinkgoT())

	//mockConfigManager = config.NewMockManager(mockCtrl)
	//externalURLs := config.ExternalURL{
	//	ProductAPI: productAPIExternalURL,
	//}
	//mockConfigManager.EXPECT().ExternalURL().Return(productAPIExternalURL).AnyTimes()

	//httpClient := httpclient.New()
	//
	//client = product.NewClient(&product.NewClientOpts{
	//	HTTPClient: httpClient,
	//	BaseURL:    productAPIExternalURL,
	//})
})

func createPact() (pact *dsl.Pact, cleanUp func()) {
	pact = &dsl.Pact{
		Host:                     "localhost",
		Consumer:                 "BasketService",
		Provider:                 "ProductService",
		DisableToolValidityCheck: true,
		PactFileWriteMode:        "merge",
		LogDir:                   "./pacts/logs",
	}

	pact.Setup(true)

	cleanUp = func() { pact.Teardown() }

	return pact, cleanUp
}
