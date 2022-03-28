package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
	"using-kit/using-kit/service"
)

const DEFAULT_PORT = ":8008"

func main() {
	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	svc := service.NewThingSvc(logger)

	getThingHandler := httptransport.NewServer(
		loggingMiddleware(log.With(logger, "method", "get-a-thing"))(service.GetAThingEndpoint(svc)),
		service.DecodeGetThingRequest,
		httptransport.EncodeJSONResponse,
	)

	getAllThingsHandler := httptransport.NewServer(
		loggingMiddleware(log.With(logger, "method", "get-aall"))(service.GetAllThings(svc)),
		service.DecodeGetAllThings,
		httptransport.EncodeJSONResponse,
	)

	r := mux.NewRouter()
	r.NotFoundHandler = http.NotFoundHandler()
	r.MethodNotAllowedHandler = http.NotFoundHandler()
	r.Handle("/things", getAllThingsHandler).Methods("GET")
	r.Handle("/thing/{id:[a-zA-Z]+}", getThingHandler).Methods("GET")

	_ = logger.Log("starting server on port:", DEFAULT_PORT)
	_ = http.ListenAndServe(DEFAULT_PORT, r)
}

func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}
