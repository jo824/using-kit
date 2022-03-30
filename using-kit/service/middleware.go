package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetAllThings(ctx context.Context) (ts []Thing, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAllThings", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAllThings(ctx)
}

func (mw loggingMiddleware) AddThing(_ context.Context, _ *Thing) error {
	panic("implement me")
}

func (mw loggingMiddleware) GetAThing(ctx context.Context, tid string) (t *Thing, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAThing", "id", tid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAThing(ctx, tid)
}
