package httpserver

import (
	"context"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/evgeny08/collection-key/types"
)

type mockService struct {
	onCreateKey       func(ctx context.Context) (*types.Key, error)
	onGetKey          func(ctx context.Context) (string, error)
	onCanceledKey     func(ctx context.Context, id string) error
	onVerificationKey func(ctx context.Context, id string) (*types.Key, error)
	onUnreleasedKey   func(ctx context.Context) ([]*types.Key, error)
}

func (s *mockService) createKey(ctx context.Context) (*types.Key, error) {
	return s.onCreateKey(ctx)
}

func (s *mockService) getKey(ctx context.Context) (string, error) {
	return s.onGetKey(ctx)
}

func (s *mockService) canceledKey(ctx context.Context, id string) error {
	return s.onCanceledKey(ctx, id)
}

func (s *mockService) verificationKey(ctx context.Context, id string) (*types.Key, error) {
	return s.onVerificationKey(ctx, id)
}

func (s *mockService) unreleasedKey(ctx context.Context) ([]*types.Key, error) {
	return s.onUnreleasedKey(ctx)
}

func startTestServer(t *testing.T) (*httptest.Server, *Client, *mockService) {
	svc := &mockService{}

	handler := newHandler(&handlerConfig{
		svc:         svc,
		logger:      log.NewNopLogger(),
		rateLimiter: rate.NewLimiter(rate.Inf, 1),
	})

	server := httptest.NewServer(handler)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	return server, client, svc
}

func TestCreateKey(t *testing.T) {
	server, client, svc := startTestServer(t)
	defer server.Close()

	testCases := []struct {
		name string
		key  *types.Key
		err  error
	}{
		{
			name: "ok response",
			key: &types.Key{
				ID:       "7777",
				Issued:   false,
				Canceled: false,
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onCreateKey = func(ctx context.Context) (*types.Key, error) {
				return tc.key, tc.err
			}
			gotKey, gotErr := client.CreateKey(context.Background())
			if !reflect.DeepEqual(gotKey, tc.key) {
				t.Fatalf("got key %#v want %#v", gotKey, tc.key)
			}
			if !reflect.DeepEqual(gotErr, tc.err) {
				t.Fatalf("got error %#v want %#v", gotErr, tc.err)
			}

		})
	}
}

func TestGetKey(t *testing.T) {
	server, client, svc := startTestServer(t)
	defer server.Close()

	testCases := []struct {
		name string
		key  *types.Key
		err  error
	}{
		{
			name: "ok response",
			key: &types.Key{
				ID:       "ki87",
				Issued:   false,
				Canceled: false,
			},
			err: nil,
		},
		{
			name: "err response",
			key: &types.Key{
				ID:       "trew",
				Issued:   true,
				Canceled: false,
			},
			err: errorf(ErrNotFound, "failed to find unreleased key"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onGetKey = func(ctx context.Context) (string, error) {
				return "", tc.err
			}
			gotKey, gotErr := client.GetKey(context.Background())
			if !reflect.DeepEqual(gotKey, tc.key) {
				t.Fatalf("got key %#v want %#v", gotKey, tc.key)
			}
			if !reflect.DeepEqual(gotErr, tc.err) {
				t.Fatalf("got error %#v want %#v", gotErr, tc.err)
			}

		})
	}
}

func TestCanceledKey(t *testing.T) {
	server, client, svc := startTestServer(t)
	defer server.Close()

	testCases := []struct {
		name string
		key  *types.Key
		err  error
	}{
		{
			name: "ok response",
			key: &types.Key{
				ID:       "ki87",
				Issued:   true,
				Canceled: false,
			},
			err: nil,
		},
		{
			name: "err response",
			key: &types.Key{
				ID:       "trew",
				Issued:   false,
				Canceled: false,
			},
			err: errorf(ErrBadParams, "failed to canceled key"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onCanceledKey = func(ctx context.Context, id string) error {
				return tc.err
			}
			gotErr := client.CanceledKey(context.Background(), tc.name)
			if !reflect.DeepEqual(gotErr, tc.err) {
				t.Fatalf("got error %#v want %#v", gotErr, tc.err)
			}

		})
	}
}

func TestVerificationKey(t *testing.T) {
	server, client, svc := startTestServer(t)
	defer server.Close()

	testCases := []struct {
		name string
		key  *types.Key
		err  error
	}{
		{
			name: "ok response",
			key: &types.Key{
				ID:       "ki87",
				Issued:   false,
				Canceled: false,
			},
			err: nil,
		},
		{
			name: "err response",
			key: &types.Key{
				ID:       "trew",
				Issued:   true,
				Canceled: false,
			},
			err: errorf(ErrNotFound, "failed to find unreleased key"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onVerificationKey = func(ctx context.Context, id string) (*types.Key, error) {
				return tc.key, tc.err
			}
			gotKey, gotErr := client.VerificationKey(context.Background(), tc.name)
			if !reflect.DeepEqual(gotKey, tc.key) {
				t.Fatalf("got key %#v want %#v", gotKey, tc.key)
			}
			if !reflect.DeepEqual(gotErr, tc.err) {
				t.Fatalf("got error %#v want %#v", gotErr, tc.err)
			}

		})
	}
}

func TestUnreleasedKey(t *testing.T) {
	server, client, svc := startTestServer(t)
	defer server.Close()

	testCases := []struct {
		name string
		key  *types.Key
		err  error
	}{
		{
			name: "ok response",
			key: &types.Key{
				ID:       "7777",
				Issued:   false,
				Canceled: false,
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onUnreleasedKey = func(ctx context.Context) ([]*types.Key, error) {
				return tc.key, tc.err
			}
			gotKey, gotErr := client.UnreleasedKey(context.Background())
			if !reflect.DeepEqual(gotKey, tc.key) {
				t.Fatalf("got key %#v want %#v", gotKey, tc.key)
			}
			if !reflect.DeepEqual(gotErr, tc.err) {
				t.Fatalf("got error %#v want %#v", gotErr, tc.err)
			}

		})
	}
}
