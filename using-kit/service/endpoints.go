package service

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetThingEndpoint     endpoint.Endpoint
	GetAllThingsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		GetThingEndpoint:     GetAThingEndpoint(svc),
		GetAllThingsEndpoint: GetAllThings(svc),
	}
}

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

func GetAllThings(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.GetAllThings(ctx)
	}
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
