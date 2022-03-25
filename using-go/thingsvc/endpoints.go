package thingsvc

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type getThingRequest struct {
	ID string `json:"id"`
}

type getThingResponse struct {
	Thing Thing  `json:"thing"`
	Err   string `json:"err,omitempty"`
}
type allThingsResponse struct {
	Ts []Thing `json:"things"`
	Err string `json:"err,omitempty"`
}
func (s *thingService) GetAllThingsEndPoint(ctx context.Context, req interface{}) (interface{}, error) {
	res, err := s.GetAllThings(ctx)
	if err != nil {
		return nil , err
	}
	return allThingsResponse{res, ""}, nil
}


func (s *thingService) GetAThingEndPoint(ctx context.Context, req interface{}) (interface{}, error) {
	r := req.(getThingRequest)
	v, err := s.Things.Find(r.ID)
	if err != nil {
		return getThingResponse{*v, err.Error()}, nil
	}
	return getThingResponse{*v, ""}, nil
}
func decodeGetThingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getThingRequest
	req.ID = mux.Vars(r)["id"]
	if len(req.ID) == 0 {
		return nil, errors.New("missing ID route param")
	}
	return req, nil
}
