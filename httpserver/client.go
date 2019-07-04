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
	createKey   endpoint.Endpoint
	getKey      endpoint.Endpoint
	canceledKey endpoint.Endpoint
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

		getKey: kithttp.NewClient(
			"GET",
			baseURL,
			encodeGetKeyRequest,
			decodeGetKeyResponse,
		).Endpoint(),

		canceledKey: kithttp.NewClient(
			"POST",
			baseURL,
			encodeCanceledKeyRequest,
			decodeCanceledKeyResponse,
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

// GetKey returns an unreleased key
func (c *Client) GetKey(ctx context.Context) (*types.Key, error) {
	var request interface{}
	response, err := c.getKey(ctx, request)
	if err != nil {
		return nil, err
	}
	res := response.(getKeyResponse)
	return res.Key, res.Err
}

// CanceledKey updates key canceled with given id
func (c *Client) CanceledKey(ctx context.Context, id string) error {
	request := canceledKeyRequest{ID: id}
	response, err := c.canceledKey(ctx, request)
	if err != nil {
		return err
	}
	res := response.(canceledKeyResponse)
	return res.Err
}
