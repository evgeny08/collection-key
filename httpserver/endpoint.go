package httpserver

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/evgeny08/collection-key/types"
)

func makeCreateKeyEndpoint(svc service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		key, err := svc.createKey(ctx)
		return createKeyResponse{Key: key, Err: err}, nil
	}
}

type createKeyResponse struct {
	Key *types.Key
	Err error
}

func makeGetKeyEndpoint(svc service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		key, err := svc.getKey(ctx)
		return getKeyResponse{Key: key, Err: err}, nil
	}
}

type getKeyResponse struct {
	Key *types.Key
	Err error
}
