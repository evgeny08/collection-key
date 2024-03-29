package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/evgeny08/collection-key/types"
)

// Service CreateKey encoders/decoders.
func encodeCreateKeyRequest(_ context.Context, r *http.Request, _ interface{}) error {
	r.URL.Path = "/api/v1/key"
	return nil
}

func decodeCreateKeyRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeCreateKeyResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(createKeyResponse)
	if res.Err != nil {
		return encodeError(w, res.Err, true)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Key)
}

func decodeCreateKeyResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return createKeyResponse{Err: decodeError(r)}, nil
	}
	res := createKeyResponse{Key: &types.Key{}}
	err := json.NewDecoder(r.Body).Decode(&res.Key)
	return res, err
}

// Service GetKey encoders/decoders.
func encodeGetKeyRequest(_ context.Context, r *http.Request, _ interface{}) error {
	r.URL.Path = "/api/v1/key/issued"
	return nil
}

func decodeGetKeyRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeGetKeyResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getKeyResponse)
	if res.Err != nil {
		return encodeError(w, res.Err, true)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&res.Key)
}

func decodeGetKeyResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return getKeyResponse{Err: decodeError(r)}, nil
	}
	res := getKeyResponse{}
	err := json.NewDecoder(r.Body).Decode(&res.Key)
	return res, err
}

// Service CanceledKey encoders/decoders.
func encodeCanceledKeyRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(canceledKeyRequest)
	r.URL.Path = "/api/v1/key/" + url.QueryEscape(req.ID) + "/canceled"
	return nil
}

func decodeCanceledKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	return canceledKeyRequest{ID: id}, nil
}

func encodeCanceledKeyResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(canceledKeyResponse)
	if res.Err != nil {
		return encodeError(w, res.Err, true)
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func decodeCanceledKeyResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return canceledKeyResponse{Err: decodeError(r)}, nil
	}
	return canceledKeyResponse{Err: nil}, nil
}

// Service VerificationKey encoders/decoders.
func encodeVerificationKeyRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(verificationKeyRequest)
	r.URL.Path = "/api/v1/key/" + url.QueryEscape(req.ID) + "/key"
	return nil
}

func decodeVerificationKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	return verificationKeyRequest{ID: id}, nil
}

func encodeVerificationKeyResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(verificationKeyResponse)
	if res.Err != nil {
		return encodeError(w, res.Err, true)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Key)
}

func decodeVerificationKeyResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return verificationKeyResponse{Err: decodeError(r)}, nil
	}
	res := verificationKeyResponse{Key: &types.Key{}}
	err := json.NewDecoder(r.Body).Decode(&res.Key)
	return res, err
}

// Service UnreleasedKey encoders/decoders.
func encodeUnreleasedKeyRequest(_ context.Context, r *http.Request, _ interface{}) error {
	r.URL.Path = "/api/v1/key"
	return nil
}

func decodeUnreleasedKeyRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeUnreleasedKeyResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(unreleasedKeyResponse)
	if res.Err != nil {
		return encodeError(w, res.Err, true)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.ListKey)
}

func decodeUnreleasedKeyResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return unreleasedKeyResponse{Err: decodeError(r)}, nil
	}
	res := unreleasedKeyResponse{ListKey: []*types.Key{}}
	err := json.NewDecoder(r.Body).Decode(&res.ListKey)
	return res, err
}

// errKindToStatus maps service error kinds to the HTTP response codes.
var errKindToStatus = map[ErrorKind]int{
	ErrBadParams: http.StatusBadRequest,
	ErrNotFound:  http.StatusNotFound,
	ErrConflict:  http.StatusConflict,
	ErrInternal:  http.StatusInternalServerError,
}

// encodeError writes a service error to the given http.ResponseWriter.
func encodeError(w http.ResponseWriter, err error, writeMessage bool) error {
	status := http.StatusInternalServerError
	message := err.Error()
	if err, ok := err.(*Error); ok {
		if s, ok := errKindToStatus[err.Kind]; ok {
			status = s
		}
		if err.Kind == ErrInternal {
			message = "internal error"
		} else {
			message = err.Message
		}
	}
	w.WriteHeader(status)
	if writeMessage {
		_, writeErr := io.WriteString(w, message)
		return writeErr
	}
	return nil
}

// decodeError reads a service error from the given *http.Response.
func decodeError(r *http.Response) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, io.LimitReader(r.Body, 1024)); err != nil {
		return fmt.Errorf("%d: %s", r.StatusCode, http.StatusText(r.StatusCode))
	}
	msg := strings.TrimSpace(buf.String())
	if msg == "" {
		msg = http.StatusText(r.StatusCode)
	}
	for kind, status := range errKindToStatus {
		if status == r.StatusCode {
			return &Error{Kind: kind, Message: msg}
		}
	}
	return fmt.Errorf("%d: %s", r.StatusCode, msg)
}
