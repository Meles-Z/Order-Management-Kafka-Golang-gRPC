package configs

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const defaultMaxBytes = 1_048_576 // 1MB

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any, maxBytes ...int) error {
	limit := defaultMaxBytes
	if len(maxBytes) > 0 {
		limit = maxBytes[0]
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(limit))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Reject unknown fields

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New("malformed JSON")
		case errors.Is(err, io.EOF):
			return errors.New("empty body")
		case errors.As(err, &unmarshalError):
			return errors.New("invalid JSON type")
		case err.Error() == "http: request body too large":
			return errors.New("request body too large")
		default:
			return err
		}
	}

	// Ensure single JSON value
	if err = dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain only one JSON value")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	w.Header().Set("Content-Type", "application/json")

	if len(headers) > 0 {
		for key, values := range headers[0] {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
