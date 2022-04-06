package service

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func BuildHTTPHandler(svc Service, l log.Logger) http.Handler {
	r := mux.NewRouter()
	eps := MakeServerEndpoints(svc)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(l)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.NotFoundHandler = http.NotFoundHandler()
	r.MethodNotAllowedHandler = http.NotFoundHandler()

	r.Methods("GET").Path("/thing/{id:[a-zA-Z]+}").Handler(httptransport.NewServer(
		eps.GetThingEndpoint,
		decodeGetThingRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/thing").Handler(httptransport.NewServer(
		eps.PostAddThing,
		decodePostThingRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/things").Handler(httptransport.NewServer(
		eps.GetAllThingsEndpoint,
		decodeGetAllThings,
		encodeResponse,
		options...,
	))

	return r
}

func decodeGetThingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getThingRequest
	req.ID = mux.Vars(r)["id"]
	if len(req.ID) == 0 {
		return nil, ErrNoID
	}
	return req, nil
}

func decodePostThingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req postThingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequestBody
	}
	return &req, nil
}

func decodeGetAllThings(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNoID:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrBadRequestBody:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
