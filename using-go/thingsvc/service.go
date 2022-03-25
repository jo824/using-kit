package thingsvc

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"using-kit/using-go/server"
)

type thingSvc interface {
	GetAThing(context.Context, string) (*Thing, error)
	GetAllThings(context.Context) ([]Thing, error)
}

type thingService struct {
	Things *ThingStore
}

func (t *thingService) Shutdown() {
	fmt.Println("goodbye")
}

func NewThingSvc() server.Service  {
	return &thingService{SeedThings()}
}

func (t *thingService) GetAThing(_ context.Context, id string) (*Thing, error) {
	return t.Things.Find(id)
}

func (t *thingService) GetAllThings(_ context.Context) ([]Thing, error) {
	return t.Things.GetAllThings()
}

func (s *thingService) Middleware(endpoint endpoint.Endpoint) endpoint.Endpoint {
	return endpoint
}

func (s *thingService) HTTPMiddleware(handler http.Handler) http.Handler {
	return handler
}

func (s *thingService) HTTPRouterOptions() []server.RouterOption {
	return nil
}

func (s *thingService) HTTPEndpoints() map[string]map[string]server.HTTPEndpoint {
	return map[string]map[string]server.HTTPEndpoint{
		"/thing/{id}": {
			"GET": {
				Endpoint: s.GetAThingEndPoint,
				Decoder:decodeGetThingRequest,
			},
		},
		"/things": {
			"GET": {
				Endpoint: s.GetAllThingsEndPoint,
			},
		},
	}
}


