package server

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
)

type Server struct {
	logger   log.Logger
	logClose func() error

	mux Router

	svc Service

	svr *http.Server

	handler http.Handler

	// exit chan for graceful shutdown
	exit chan chan error
}

func NewServer(svc Service) *Server {
	// load config from environment with defaults set
	ropts := svc.HTTPRouterOptions()
	// default the router if none set
	if len(ropts) == 0 {
		ropts = append(ropts, RouterSelect(""))
	}
	var r Router
	for _, opt := range ropts {
		r = opt(r)
	}

	//ctx := context.Background()

	s := &Server{
		logger: newJSONLogger(),
		logClose: func() error { return nil },
		mux:  r,
		exit: make(chan chan error),
	}
	s.svr = &http.Server{
		Handler: s.mux,
		Addr:    ":8833",
	}
	s.register(svc)
	return s
}

func (s *Server) register(svc Service) {
	s.svc = svc
	s.handler = s.svc.HTTPMiddleware(s.mux)

	// register all endpoints with our wrappers & default decoders/encoders
	for path, epMethods := range svc.HTTPEndpoints() {
		for method, ep := range epMethods {
			// just pass the http.Request in if no decoder provided
			if ep.Decoder == nil {
				ep.Decoder = basicDecoder
			}
			// default to the httptransport helper
			if ep.Encoder == nil {
				ep.Encoder = httptransport.EncodeJSONResponse
			}
			s.mux.Handle(method, path,
				httptransport.NewServer(
					svc.Middleware(ep.Endpoint),
					ep.Decoder,
					ep.Encoder))
		}
	}

}
func basicDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}
func (s *Server) start() error {
	go func() {
		err := s.svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Log(
				"error", err,
				"message", "HTTP server error - initiating shutting down")
			s.stop()
		}
	}()

	s.logger.Log("message",
		fmt.Sprintf("listening on HTTP port: %s", ":8833"))


	go func() {
		exit := <-s.exit

		// stop the listener with timeout
		dur, _ := time.ParseDuration("30s")
		ctx, cancel := context.WithTimeout(context.Background(), dur)
		defer cancel()
		defer func() {
			// flush the logger after server shuts down
			if s.logClose != nil {
				s.logClose()
			}

		}()

		if shutdown, ok := s.svc.(Shutdowner); ok {
			shutdown.Shutdown()
		}

		exit <- s.svr.Shutdown(ctx)
	}()

	return nil
}

func (s *Server) stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

// Run Service and start up the server(s).
// This will block until the server shuts down.
func Run(service Service) error {
	svr := NewServer(service)

	if err := svr.start(); err != nil {
		return err
	}

	// parse address for host, port
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	svr.logger.Log("received signal", <-ch)
	return svr.stop()
}

func newJSONLogger() log.Logger {
	return log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
}
