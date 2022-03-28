package service

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
)

func GetAThingEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getThingRequest)
		v, err := svc.GetAThing(ctx, req.ID)
		if err != nil && v == nil {
			return getThingResponse{}, errors.New("requested thing doesn't exist\n")
		}
		return getThingResponse{*v, ""}, nil
	}
}

func DecodeGetThingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getThingRequest
	req.ID = mux.Vars(r)["id"]
	if len(req.ID) == 0 {
		return nil, errors.New("missing ID route param")
	}
	return req, nil
}

func GetAllThings(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.GetAllThings(ctx)
	}
}

func DecodeGetAllThings(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

type getThingRequest struct {
	ID string `json:"id"`
}

type getThingResponse struct {
	Thing Thing  `json:"thing"`
	Err   string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type AllThingsResponse struct {
	Ts []Thing `json:"v"`
}
