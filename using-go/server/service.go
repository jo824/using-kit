package server

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

// HTTPEndpoint encapsulates everything required to build
// an endpoint hosted on a kit server.
type HTTPEndpoint struct {
	Endpoint endpoint.Endpoint
	Decoder  httptransport.DecodeRequestFunc
	Encoder  httptransport.EncodeResponseFunc
}

type Service interface {
	HTTPEndpoints() map[string]map[string]HTTPEndpoint
	Middleware(endpoint.Endpoint) endpoint.Endpoint
	HTTPMiddleware(http.Handler) http.Handler
	HTTPRouterOptions() []RouterOption
}

type Shutdowner interface {
	Shutdown()
}