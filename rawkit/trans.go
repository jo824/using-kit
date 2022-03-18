package rawkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
)

//Now we need to expose your service to the outside world, so it can be called.
//Your organization probably already has opinions about how services should talk to each other.
//Maybe you use Thrift, or custom JSON over HTTP. Go kit supports many transports out of the box.
//For this minimal example service, let’s use JSON over HTTP. Go kit provides a helper struct, in package transport/http.

func DecodeGetThingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getThingRequest
	req.ID = mux.Vars(r)["id"]
	if len(req.ID) == 0 {
		return nil, errors.New("missing ID route param")
	}
	return req, nil
}

func DecodeGetAllThings(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

//Go kit provides much of its functionality through an abstraction called an endpoint.
//An endpoint is defined as follows (you don’t have to put it anywhere in the code, it is provided by go-kit):
// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

//It represents a single RPC.
//That is, a single method in our service interface.
//We’ll write simple adapters to convert each of our service’s methods into an endpoint.
//Each adapter takes a ThingService, and returns an endpoint that corresponds to one of the methods.

func GetAThingEP(svc StrSvc) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getThingRequest)
		// propagate context later for tracing
		v, err := svc.GetAThing(context.TODO(), req.ID)
		fmt.Printf("the thing:%+v", v)
		if err != nil {
			return getThingResponse{*v, err.Error()}, nil
		}
		return getThingResponse{*v, ""}, nil
	}
}

func GetAllThings(svc StrSvc) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return svc.GetAllThings(context.TODO())
	}
}
