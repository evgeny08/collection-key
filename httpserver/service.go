package httpserver

import (
	"context"
	"math/rand"
	"strings"

	"github.com/go-kit/kit/log"

	"github.com/evgeny08/collection-key/types"
)

// service manages HTTP server methods.
type service interface {
	createKey(ctx context.Context) (*types.Key, error)
	getKey(ctx context.Context) (*types.Key, error)
	canceledKey(ctx context.Context, id string) error
}

type basicService struct {
	logger  log.Logger
	storage Storage
}

// createKey creates a new key
func (s *basicService) createKey(ctx context.Context) (*types.Key, error) {
	keyLength := 4
	key := &types.Key{
		ID:       genKey(keyLength),
		Issued:   false,
		Canceled: false,
	}
	err := s.storage.InsertKey(ctx, key)
	if err != nil {
		return nil, errorf(ErrBadParams, "failed to insert key: %v", err)
	}
	return key, nil
}

// genKey generates a key with given length
func genKey(n int) string {
	letterRunes := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GetKey returns an unreleased key
func (s *basicService) getKey(ctx context.Context) (*types.Key, error){
	key, err := s.storage.GetKey(ctx)
	if err != nil {
		return nil, errorf(ErrBadParams, "failed to find unreleased key: %v", err)
	}
	return key, nil
}

// canceledKey updates key canceled with given id
func (s *basicService) canceledKey(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errorf(ErrBadParams, "empty key id")
	}

	err := s.storage.CanceledKey(ctx, id)
	if err != nil {
		if storageErrIsNotFound(err) {
			return errorf(ErrNotFound, "key is not found")
		}
		return errorf(ErrBadParams, "failed to canceled key: %v", err)
	}
	return nil
}

// storageErrIsNotFound checks if the storage error is "not found".
func storageErrIsNotFound(err error) bool {
	type notFound interface {
		NotFound() bool
	}
	e, ok := err.(notFound)
	return ok && e.NotFound()
}
