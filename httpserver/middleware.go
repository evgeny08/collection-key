package httpserver

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/evgeny08/collection-key/types"
)

// loggingMiddleware wraps Service and logs request information to the provided logger.
type loggingMiddleware struct {
	next   service
	logger log.Logger
}

func (m *loggingMiddleware) createKey(ctx context.Context) (*types.Key, error) {
	begin := time.Now()
	key, err := m.next.createKey(ctx)
	err = level.Info(m.logger).Log(
		"method", "CreateKey",
		"err", err,
		"elapsed", time.Since(begin),
		"id", key.ID,
	)
	return key, err
}

func (m *loggingMiddleware) getKey(ctx context.Context) (*types.Key, error) {
	begin := time.Now()
	key, err := m.next.getKey(ctx)
	err = level.Info(m.logger).Log(
		"method", "GetKey",
		"err", err,
		"elapsed", time.Since(begin),
		"id", key.ID,
	)
	return key, err
}

func (m *loggingMiddleware) canceledKey(ctx context.Context, id string) error {
	begin := time.Now()
	err := m.next.canceledKey(ctx, id)
	err = level.Info(m.logger).Log(
		"method", "CanceledKey",
		"err", err,
		"elapsed", time.Since(begin),
		"id", id,
	)
	return err
}
