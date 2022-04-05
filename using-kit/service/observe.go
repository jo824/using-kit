package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type observeService struct {
	count   metrics.Counter
	latency metrics.Histogram
	Service
}

func NewObserveService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &observeService{
		count:   counter,
		latency: latency,
		Service: s,
	}
}

func (s *observeService) GetAThing(ctx context.Context, id string) (*Thing, error) {
	defer func(begin time.Time) {
		s.count.With("method", "GetAThing").Add(1)
		s.latency.With("method", "GetAThing").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetAThing(ctx, id)
}

func (s *observeService) GetAllThings(ctx context.Context) ([]Thing, error) {
	defer func(begin time.Time) {
		s.count.With("method", "GetAllThings").Add(1)
		s.latency.With("method", "GetAllThings").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetAllThings(ctx)
}

func (s *observeService) AddThing(_ context.Context, _ *Thing) error {
	panic("implement me")
}
