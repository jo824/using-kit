package service

import (
	"context"
	"errors"
)

var (
	ErrNoID           = errors.New("id doesn't exist")
	ErrAlreadyExists  = errors.New("thing with id already exists")
	ErrBadRequestBody = errors.New("bad request body")
)

type Service interface {
	GetAThing(context.Context, string) (*Thing, error)
	GetAllThings(context.Context) ([]Thing, error)
	AddThing(ctx context.Context, thing *Thing) error
}

type ThingService struct {
	Things *ThingStore
}

func NewThingSvc() Service {
	return ThingService{SeedThings()}
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
