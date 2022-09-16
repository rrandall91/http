package http

import (
	"io"
	"time"
)

// Response represents an HTTP response.
type Response struct {
	StatusCode int
	Duration   time.Duration
	Body       io.Reader
	Headers    []Param
}
