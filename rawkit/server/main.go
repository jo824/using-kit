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
	"using-kit/rawkit"
)

const DEFAULT_PORT = ":8080"

func main() {
	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	svc := rawkit.NewThingSvc(logger)

	getThingHandler := httptransport.NewServer(
		loggingMiddleware(log.With(logger, "method", "get-a-thing"))(rawkit.GetAThingEP(svc)),
		rawkit.DecodeGetThingRequest,
		httptransport.EncodeJSONResponse,
	)

	getAllThings := httptransport.NewServer(
		loggingMiddleware(log.With(logger, "method", "get-aall"))(rawkit.GetAllThings(svc)),
		rawkit.DecodeGetAllThings,
		httptransport.EncodeJSONResponse,
	)

	r := mux.NewRouter()
	r.NotFoundHandler = http.NotFoundHandler()
	r.MethodNotAllowedHandler = http.NotFoundHandler()
	r.Handle("/things", getAllThings).Methods("GET")
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
