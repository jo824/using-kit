package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Service interface {
	GetAThing(context.Context, string) (*Thing, error)
	GetAllThings(context.Context) ([]Thing, error)
	AddThing(ctx context.Context, thing *Thing) error
}

type ThingService struct {
	Things *ThingStore
	l      log.Logger
}

func NewThingSvc(l log.Logger) Service {
	return ThingService{SeedThings(), l}
}

func (t ThingService) GetAThing(_ context.Context, id string) (*Thing, error) {
	return t.Things.Find(id)
}

func (t ThingService) GetAllThings(_ context.Context) ([]Thing, error) {
	return t.Things.GetAllThings()
}

func (t ThingService) AddThing(_ context.Context, thing *Thing) error {
	return t.Things.Save(thing)
}