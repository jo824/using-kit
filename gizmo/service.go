package gizzz

import (
	"context"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"using-kit/shared"

	httptransport "github.com/go-kit/kit/transport/http"
)

type service struct {
	lg     log.Logger
	things *shared.ThingStore
}

func NewService() kit.Service {
	l := log.NewLogfmtLogger(os.Stderr)
	l = level.NewFilter(l, level.AllowAll())
	l = log.With(l, "ts", log.DefaultTimestamp)

	store := shared.SeedThings()
	return &service{l, store}
}

func (s *service) HTTPOptions() []httptransport.ServerOption {
	return []httptransport.ServerOption{}
}

// There are multiple router options the deafult is
// https://github.com/gorilla/mux
// github.com/NYTimes/gizmo/server/kit/router.go
func (s *service) HTTPRouterOptions() []kit.RouterOption {
	return []kit.RouterOption{
		kit.RouterSelect("gorilla"),
	}
}

func (s *service) HTTPMiddleware(h http.Handler) http.Handler {
	return h
}

func (s *service) Middleware(ep endpoint.Endpoint) endpoint.Endpoint {
	return endpoint.Endpoint(func(ctx context.Context, r interface{}) (interface{}, error) {
		rr := r.(*http.Request)
		path := rr.URL.Path
		_ = s.lg.Log("route:", path)
		return ep(ctx, r)
	})
}

//Declare available endpoints
func (s *service) HTTPEndpoints() map[string]map[string]kit.HTTPEndpoint {
	return map[string]map[string]kit.HTTPEndpoint{
		"/things": {
			"GET": {
				Endpoint: s.getThings,
			},
		},
		"/things/{id:[a-zA-Z]+}": {
			"GET": {
				Endpoint: s.getAThing,
				Decoder:  decodeGetRequest,
			},
		},
	}
}

// Just satisfiying the interface with next 3
func (s *service) RPCMiddleware() grpc.UnaryServerInterceptor {
	return nil
}

func (s *service) RPCServiceDesc() *grpc.ServiceDesc {
	return nil
}

func (s *service) RPCOptions() []grpc.ServerOption {
	return nil
}

// go-kit endpoint.Endpoint with core business logic
func (s *service) getThings(_ context.Context, _ interface{}) (interface{}, error) {
	return s.things.GetAllThings()
}

func (s *service) getAThing(_ context.Context, req interface{}) (interface{}, error) {
	id := req.(*GetThingsReq).ID
	s.lg.Log("looking for id:", id)
	return s.things.Find(id)
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tid, _ := kit.Vars(r)["id"]
	return &GetThingsReq{
		ID: tid,
	}, nil
}

type Message struct {
	Message string `json:"message,omitempty"`
}

type GetThingsReq struct {
	ID string `json:"id,omitempty"`
}
