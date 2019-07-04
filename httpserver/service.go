package httpserver

import (
	"context"
	"math/rand"

	"github.com/go-kit/kit/log"

	"github.com/evgeny08/collection-key/types"
)

// service manages HTTP server methods.
type service interface {
	createKey(ctx context.Context) (*types.Key, error)
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

// storageErrIsNotFound checks if the storage error is "not found".
func storageErrIsNotFound(err error) bool {
	type notFound interface {
		NotFound() bool
	}
	e, ok := err.(notFound)
	return ok && e.NotFound()
}
