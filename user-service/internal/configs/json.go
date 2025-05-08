package configs

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

const defaultMaxBytes = 1_048_576 // 1MB

// ReadJSON reads JSON from the request body into the provided destination.
// It handles body size limits and common JSON decoding errors.
func ReadJSON(c echo.Context, dst any, maxBytes ...int) error {
	req := c.Request()
	res := c.Response()

	limit := defaultMaxBytes
	if len(maxBytes) > 0 {
		limit = maxBytes[0]
	}

	req.Body = http.MaxBytesReader(res, req.Body, int64(limit))

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Reject unknown fields

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return echo.NewHTTPError(http.StatusBadRequest, "malformed JSON")
		case errors.Is(err, io.EOF):
			return echo.NewHTTPError(http.StatusBadRequest, "empty body")
		case errors.As(err, &unmarshalError):
			return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON type")
		case err.Error() == "http: request body too large":
			return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "request body too large")
		default:
			return err
		}
	}

	// Ensure single JSON value
	if err = dec.Decode(&struct{}{}); err != io.EOF {
		return echo.NewHTTPError(http.StatusBadRequest, "body must contain only one JSON value")
	}

	return nil
}

// WriteJSON writes JSON data to the response with the specified status code and headers.
func WriteJSON(c echo.Context, status int, data any, headers ...http.Header) error {
	res := c.Response()
	res.Header().Set("Content-Type", "application/json")

	if len(headers) > 0 {
		for key, values := range headers[0] {
			for _, value := range values {
				res.Header().Add(key, value)
			}
		}
	}

	res.WriteHeader(status)
	return json.NewEncoder(res).Encode(data)
}
