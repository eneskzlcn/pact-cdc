package product_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/eneskzlcn/pact-cdc/provider/product-service/app/product"
	"github.com/eneskzlcn/pact-cdc/server"
	"github.com/golang/mock/gomock"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
	"time"
)

const (
	testDBUser         = "test"
	testDBPass         = "test"
	testDBName         = "test"
	pactBrokerLocalURL = "http://localhost"
)

type PactSettings struct {
	Host            string
	ProviderName    string
	BrokerBaseURL   string
	BrokerUsername  string // Basic authentication
	BrokerPassword  string // Basic authentication
	ConsumerName    string
	ConsumerVersion string // a git sha, semantic version number
	ConsumerTag     string // dev, staging, prod
	ProviderVersion string
}

func (s *PactSettings) getPactURL(useLocal bool) string {
	var pactURL string

	if s.ConsumerVersion == "" {
		pactURL = fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/latest/master.json", s.BrokerBaseURL, s.ProviderName, s.ConsumerName)
	} else {
		pactURL = fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/version/%s.json", s.BrokerBaseURL, s.ProviderName, s.ConsumerName, s.ConsumerVersion)
	}

	return pactURL
}

type ProviderTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	pactSettings *PactSettings
	ctx          context.Context
	l            *logrus.Logger
	app          server.Server
	mockRepo     *product.MockRepository
	serverPort   string
}

func TestProvider(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}

func (s *ProviderTestSuite) SetupSuite() {
	s.l, _ = test.NewNullLogger()
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = product.NewMockRepository(s.ctrl)
	//req := testcontainers.ContainerRequest{
	//	Image:        "postgres:latest",
	//	ExposedPorts: []string{"5432/tcp"},
	//	Env: map[string]string{
	//		"POSTGRES_DB":       testDBName,
	//		"POSTGRES_USER":     testDBUser,
	//		"POSTGRES_PASSWORD": testDBPass,
	//	},
	//	WaitingFor: wait.ForListeningPort("5432/tcp"),
	//}

	//postgresContainer, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
	//	ContainerRequest: req,
	//	Started:          true,
	//})
	//s.Nil(err)
	//
	//host, _ := postgresContainer.Host(s.ctx)
	//p, _ := postgresContainer.MappedPort(s.ctx, "5432/tcp")
	//port := p.Port()
	//
	//postgreDB := postgres.New(&postgres.NewPostgresOpts{
	//	Host:     host,
	//	Port:     port,
	//	Username: testDBUser,
	//	Password: testDBPass,
	//	DBName:   testDBName,
	//})
	//
	//postgresRepository := persistence.NewPostgresRepository(&persistence.NewPostgresRepositoryOpts{
	//	DB: postgreDB,
	//	L:  s.l,
	//})

	productService := product.NewService(&product.NewServiceOpts{
		R: s.mockRepo,
		L: s.l,
	})

	productHandler := product.NewHandler(&product.NewHandlerOpts{
		S: productService,
		L: s.l,
	})

	sp, err := utils.GetFreePort()
	s.Nil(err)

	s.serverPort = fmt.Sprintf("%d", sp)

	s.app = server.New(&server.NewServerOpts{
		Port: s.serverPort,
	}, []server.RouteHandler{
		productHandler,
	})

	//err = createProductTableOnDB(postgreDB)
	s.Nil(err)

	go func() {
		if serverErr := s.app.Run(); serverErr != nil {
			fmt.Println("serverErr", serverErr)
		}
	}()

	_ = os.Setenv("CONSUMER_NAME", "BasketService")
	_ = os.Setenv("CONSUMER_TAG", "dev")
	_ = os.Setenv("GIT_SHORT_HASH", "4.0.2")
	_ = os.Setenv("CONSUMER_VERSION", "4.0.2")
	s.pactSettings = &PactSettings{
		Host:            "localhost",
		ProviderName:    "ProductService",
		ConsumerName:    os.Getenv("CONSUMER_NAME"),
		ConsumerVersion: os.Getenv("CONSUMER_VERSION"),
		BrokerBaseURL:   pactBrokerLocalURL,
		ConsumerTag:     os.Getenv("CONSUMER_TAG"),
		ProviderVersion: os.Getenv("GIT_SHORT_HASH"),
	}
	time.Sleep(3 * time.Second)
}

func (s *ProviderTestSuite) TestProvider() {
	pact := &dsl.Pact{
		Host:                     s.pactSettings.Host,
		Provider:                 s.pactSettings.ProviderName,
		Consumer:                 s.pactSettings.ConsumerName,
		DisableToolValidityCheck: true,
	}

	providerBaseURL := fmt.Sprintf("http://%s:%s", s.pactSettings.Host, s.serverPort)

	verifyRequest := types.VerifyRequest{
		ProviderBaseURL:            providerBaseURL,
		PactURLs:                   []string{s.pactSettings.getPactURL(true)},
		BrokerURL:                  s.pactSettings.BrokerBaseURL,
		Tags:                       []string{s.pactSettings.ConsumerTag},
		BrokerUsername:             s.pactSettings.BrokerUsername,
		BrokerPassword:             s.pactSettings.BrokerPassword,
		FailIfNoPactsFound:         true,
		PublishVerificationResults: true,
		ProviderVersion:            s.pactSettings.ProviderVersion,
		StateHandlers: map[string]types.StateHandler{
			//  /products/bulk endpoints provider states
			/*	"i get validation error when no product id is given":                                  s.iGetValidationErrorWhenNoProductIDIsGivenStateHandler,
				"i get product not found error when the one of product with given id does not exists": s.iGetProductNotFoundErrorWhenTheOneOfProductWithGivenIDDoesNotExistsStateHandler,
				"i get products with given ids":                                                       s.iGetProductsWithGivenIDsStateHandler,*/

			// /products/{id} endpoints provider states
			"i get product with given id": s.iGetProductWithGivenIDStateHandler,
			"i get product not found error when the product with given id does not exists": s.iGetProductNotFoundErrorWhenTheProductWithGivenIDDoesNotExistsStateHandler,
		},
		BeforeEach:       nil,
		AfterEach:        nil,
		TagWithGitBranch: false,
	}
	defer pact.Teardown()
	verifyResponses, err := pact.VerifyProvider(s.T(), verifyRequest)
	s.Nil(err)

	if err != nil {
		log.Println(err)
	}

	log.Printf("%d pact tests run", len(verifyResponses))
}

func (s *ProviderTestSuite) iGetValidationErrorWhenNoProductIDIsGivenStateHandler() error {
	//no need to do anything, automatically captured on handler layer.
	return nil
}

func (s *ProviderTestSuite) iGetProductWithGivenIDStateHandler() error {
	/*
		id = "52fdfc07-2182-454f-963f-5f0f9a621d72"

		Expected ID In Contract.
	*/

	givenProductID := "52fdfc07-2182-454f-963f-5f0f9a621d72"

	s.mockRepo.EXPECT().GetProductByID(gomock.Any(), givenProductID).
		Return(randomProduct(givenProductID), nil)

	return nil
}

func (s *ProviderTestSuite) iGetProductNotFoundErrorWhenTheProductWithGivenIDDoesNotExistsStateHandler() error {
	/*
		id = "52fdfc07-2182-454f-963f-5f0f9a621d72"

		Expected ID In Contract.
	*/
	givenProductID := "52fdfc07-2182-454f-963f-5f0f9a621d72"
	var _ = givenProductID

	s.mockRepo.EXPECT().GetProductByID(gomock.Any(), givenProductID).Return(nil, sql.ErrNoRows)
	return nil
}

func (s *ProviderTestSuite) iGetProductNotFoundErrorWhenTheOneOfProductWithGivenIDDoesNotExistsStateHandler() error {
	/*
			"ids": [
				"9566c74d-1003-4c4d-bbbb-0407d1e2c649",
				"81855ad8-681d-4d86-91e9-1e00167939cb",
				"6694d2c4-22ac-4208-a007-2939487f6999"
			]

		ids expected in contract and we will not give one of them to see if the state works as expected.
	*/
	//notExistID := "81855ad8-681d-4d86-91e9-1e00167939cb"
	//
	//products := []*product.Product{
	//	randomProduct("9566c74d-1003-4c4d-bbbb-0407d1e2c649"),
	//	randomProduct("6694d2c4-22ac-4208-a007-2939487f6999"),
	//}
	//var _ = notExistID
	//s.mockRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(products[0], nil).Times(1)
	//s.mockRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(products[1], nil).Times(1)
	//s.mockRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows).Times(1)

	return nil
}

func (s *ProviderTestSuite) iGetProductsWithGivenIDsStateHandler() error {
	/*
			"ids": [
				"9566c74d-1003-4c4d-bbbb-0407d1e2c649",
				"81855ad8-681d-4d86-91e9-1e00167939cb",
				"6694d2c4-22ac-4208-a007-2939487f6999"
			]

		ids expected in contract and we will give all of them to see if the state works as expected.
	*/
	//products := []*product.Product{
	//	randomProduct("9566c74d-1003-4c4d-bbbb-0407d1e2c649"),
	//	randomProduct("6694d2c4-22ac-4208-a007-2939487f6999"),
	//	randomProduct("81855ad8-681d-4d86-91e9-1e00167939cb"),
	//}
	//s.mockRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(products[0], nil).Times(1)
	//s.mockRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(products[1], nil).Times(1)
	//s.mockRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(products[2], nil).Times(1)

	return nil
}

/*
func createProductTableOnDB(sql *sql.DB) error {
	_, err := sql.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(255) NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price NUMERIC(10,2) NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    buying_price NUMERIC(10,2) NOT NULL,
		    selling_price NUMERIC(10,2) NOT NULL,
		    image_url VARCHAR(255) NOT NULL,
		    type VARCHAR(255) NOT NULL,
		    provider VARCHAR(255) NOT NULL,
		    creator VARCHAR(255) NOT NULL,
		    distributor VARCHAR(255) NOT NULL,
		    code VARCHAR(255) NOT NULL,
		    color VARCHAR(255) NOT NULL
		);
	`)

	return err
}
*/

func randomProduct(id string) *product.Product {
	if id == "" {
		id = gofakeit.UUID()
	}
	return &product.Product{
		ID:           id,
		Name:         gofakeit.Name(),
		Code:         gofakeit.Word(),
		Color:        gofakeit.Color(),
		CreatedAt:    gofakeit.Date(),
		UpdatedAt:    gofakeit.Date(),
		BuyingPrice:  gofakeit.Price(0, 3000),
		SellingPrice: gofakeit.Price(3500, 10000),
		ImageURL:     gofakeit.ImageURL(100, 200),
		Type: product.ProductType(
			gofakeit.RandString([]string{
				string(product.Bag), string(product.Hat), string(product.Clothing),
			})),
		Provider:    gofakeit.Company(),
		Creator:     gofakeit.Company(),
		Distributor: gofakeit.Company(),
	}
}
