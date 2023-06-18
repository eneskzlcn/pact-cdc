package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/eneskzlcn/pact-cdc/consumer/basket-service/app/cerr"
)

type Client interface {
	Get(ctx context.Context, url string, headers map[string]string) ([]byte, error)
	GetWithBody(
		ctx context.Context,
		url string,
		headers map[string]string,
		body interface{},
	) ([]byte, error)
}

type client struct {
	httpClient *http.Client
}

func New() Client {
	return &client{
		httpClient: http.DefaultClient,
	}
}

func (c *client) Get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var bag cerr.Bag
		if err := json.NewDecoder(resp.Body).Decode(&bag); err != nil {
			return nil, err
		}

		return nil, bag
	}

	return io.ReadAll(resp.Body)

}

func (c *client) GetWithBody(
	ctx context.Context, url string, headers map[string]string, body interface{}) ([]byte, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var bag cerr.Bag
		if err := json.NewDecoder(resp.Body).Decode(&bag); err != nil {
			return nil, err
		}

		return nil, bag
	}

	return io.ReadAll(resp.Body)

}

var DefaultHeaders = map[string]string{
	"Content-Type": "application/json",
	"Acccept":      "*/*",
}
