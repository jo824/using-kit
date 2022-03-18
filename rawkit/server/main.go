package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"using-kit/rawkit"
	"using-kit/shared"
)

// Any component that needs to log should treat the logger like a dependency, same as a database connection.
// So, we construct our logger in our func main, and pass it to components that need it.
// We never use a globally-scoped logger.
// We could pass a logger directly into our stringService implementation, but there’s a better way.
// Let’s use a middleware, also known as a decorator. A middleware is a function that takes an endpoint
// and returns an endpoint.
// type Middleware func(Endpoint) Endpoint
// Note, that the Middleware type is provided for you by go-kit.
const DEFAULT_PORT = ":8080"

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	svc := rawkit.ThingService{Things: shared.SeedThings()}

	getThingHandler := httptransport.NewServer(
		loggingMiddleware(log.With(logger, "method", "get-a-thing"))(rawkit.GetAThingEP(svc)),
		rawkit.DecodeGetThingRequest,
		rawkit.EncodeResponse,
	)

	getAllThings := httptransport.NewServer(
		loggingMiddleware(log.With(logger, "method", "get-aall"))(rawkit.GetAllThings(svc)),
		rawkit.DecodeGetAllThings,
		rawkit.EncodeResponse,
	)

	r := mux.NewRouter()
	r.Handle("/things", getAllThings)
	r.Handle("/things/{id:[a-zA-Z]+}", getThingHandler)

	logger.Log(`starting server on port:`, DEFAULT_PORT)
	//returns err
	_ = http.ListenAndServe(DEFAULT_PORT, r)
}

func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "calling endpoint")
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}
