package httpserver

import (
	"context"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/evgeny08/collection-key/types"
)

// Client is a client for auth-user service.
type Client struct {
	createKey endpoint.Endpoint
}

// NewClient creates a new service client.
func NewClient(serviceURL string) (*Client, error) {
	baseURL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		createKey: kithttp.NewClient(
			"POST",
			baseURL,
			encodeCreateKeyRequest,
			decodeCreateKeyResponse,
		).Endpoint(),
	}

	return c, nil
}

// CreateKey creates a new key.
func (c *Client) CreateKey(ctx context.Context) (*types.Key, error) {
	var request interface{}
	response, err := c.createKey(ctx, request)
	if err != nil {
		return nil, err
	}
	res := response.(createKeyResponse)
	return res.Key, res.Err
}
