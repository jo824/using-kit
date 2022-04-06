package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetThingEndpoint     endpoint.Endpoint
	GetAllThingsEndpoint endpoint.Endpoint
	PostAddThing         endpoint.Endpoint
}

func MakeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		GetThingEndpoint:     GetAThingEndpoint(svc),
		GetAllThingsEndpoint: GetAllThings(svc),
		PostAddThing:         PostAddThing(svc),
	}
}

func GetAThingEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getThingRequest)
		v, err := svc.GetAThing(ctx, req.ID)
		if err != nil && v == nil {
			return nil, err
		}
		return getThingResponse{*v, ""}, nil
	}
}

func GetAllThings(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.GetAllThings(ctx)
	}
}

func PostAddThing(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*postThingRequest)
		t := &Thing{
			ID:        req.ID,
			Available: req.Available,
		}
		err := svc.AddThing(ctx, t)
		if err != nil {
			return nil, err
		}
		return postThingResponse{"successfully added"}, nil
	}

}

type getThingRequest struct {
	ID string `json:"id"`
}

type getThingResponse struct {
	Thing Thing  `json:"thing"`
	Err   string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type postThingRequest struct {
	ID        string `json:"id"`
	Available bool   `json:"available"`
}
type postThingResponse struct {
	Message string `json:"message"`
}

type allThingsResponse struct {
	Ts []Thing `json:"v"`
}
