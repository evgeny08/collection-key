package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/evgeny08/collection-key/types"
)

// Service.CreateKey encoders/decoders.
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
