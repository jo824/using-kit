package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"using-kit/service"
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

	svc := service.NewThingSvc()
	svc = service.LoggingMiddleware(l)(svc)
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
