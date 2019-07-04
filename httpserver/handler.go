package httpserver

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

type handlerConfig struct {
	svc         service
	logger      log.Logger
	rateLimiter *rate.Limiter
}

// newHandler creates a new HTTP handler serving service endpoints.
func newHandler(cfg *handlerConfig) http.Handler {
	svc := &loggingMiddleware{next: cfg.svc, logger: cfg.logger}

	createKeyEndpoint := makeCreateKeyEndpoint(svc)
	createKeyEndpoint = applyMiddleware(createKeyEndpoint, "CreateKey", cfg)

	getKeyEndpoint := makeGetKeyEndpoint(svc)
	getKeyEndpoint = applyMiddleware(getKeyEndpoint, "GetKey", cfg)

	canceledKeyEndpoint := makeCanceledKeyEndpoint(svc)
	canceledKeyEndpoint = applyMiddleware(canceledKeyEndpoint, "RedemptionKey", cfg)

	router := mux.NewRouter()

	router.Path("/api/v1/key").Methods("POST").Handler(kithttp.NewServer(
		createKeyEndpoint,
		decodeCreateKeyRequest,
		encodeCreateKeyResponse,
	))

	router.Path("/api/v1/key/issued").Methods("GET").Handler(kithttp.NewServer(
		getKeyEndpoint,
		decodeGetKeyRequest,
		encodeGetKeyResponse,
	))

	router.Path("/api/v1/key/{id}/canceled").Methods("POST").Handler(kithttp.NewServer(
		canceledKeyEndpoint,
		decodeCanceledKeyRequest,
		encodeCanceledKeyResponse,
	))

	return router
}

func applyMiddleware(e endpoint.Endpoint, method string, cfg *handlerConfig) endpoint.Endpoint {
	return ratelimit.NewErroringLimiter(cfg.rateLimiter)(e)
}
