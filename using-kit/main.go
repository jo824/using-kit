package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"using-kit/using-kit/service"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const DEFAULT_PORT = ":8008"

func main() {
	// Create a single logger, which we'll use and give to other components.
	var l log.Logger
	{
		l = log.NewLogfmtLogger(os.Stderr)
		l = log.With(l, "ts", log.DefaultTimestampUTC)
		l = log.With(l, "caller", log.DefaultCaller)
	}
	obsKeys := []string{"method"}

	svc := service.NewThingSvc()
	svc = service.LoggingMiddleware(l)(svc)
	svc = service.NewObserveService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "thing_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, obsKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "thing_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, obsKeys),
		svc,
	)

	r := service.BuildHTTPHandler(svc, log.With(l, "component", "HTTP"))

	//Here we are taking advantage of go routines and channels
	//so that when the server shuts down it can do it gracefully
	//logging any errors or taking advantage of exit functions provided by the server impl
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		l.Log("transport", "HTTP", "addr", DEFAULT_PORT)
		errs <- http.ListenAndServe(DEFAULT_PORT, r)
	}()

	l.Log("exit", <-errs)
}
