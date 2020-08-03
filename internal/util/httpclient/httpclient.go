package httpclient

import (
	"net/http"
)

// Factory func
type Factory func() *http.Client

// DefaultFactory func
func DefaultFactory() *http.Client {
	return &http.Client{}
}

// RoundTripFunc type
type RoundTripFunc func(*http.Request) *http.Response

// RoundTrip func
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewFactoryMock func
func NewFactoryMock(fn RoundTripFunc) Factory {
	return func() *http.Client {
		return &http.Client{
			Transport: RoundTripFunc(fn),
		}
	}
}
