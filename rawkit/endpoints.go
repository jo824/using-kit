package rawkit

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
	"using-kit/shared"
)

func GetAThingEP(svc ThingSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetThingRequest)
		v, err := svc.GetAThing(ctx, req.ID)
		if err != nil {
			return GetThingResponse{*v, err.Error()}, nil
		}
		return GetThingResponse{*v, ""}, nil
	}
}

func GetAllThings(svc ThingSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.GetAllThings(ctx)
	}
}

func DecodeGetThingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetThingRequest
	req.ID = mux.Vars(r)["id"]
	if len(req.ID) == 0 {
		return nil, errors.New("missing ID route param")
	}
	return req, nil
}

func DecodeGetAllThings(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

type GetThingRequest struct {
	ID string `json:"id"`
}

type GetThingResponse struct {
	Thing shared.Thing `json:"thing"`
	Err   string       `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type AllThingsResponse struct {
	Ts []shared.Thing `json:"v"`
}
