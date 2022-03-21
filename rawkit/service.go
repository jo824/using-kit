package rawkit

import (
	"context"
	"github.com/go-kit/kit/log"
	"using-kit/shared"
)

type ThingSvc interface {
	GetAThing(context.Context, string) (*shared.Thing, error)
	GetAllThings(context.Context) ([]shared.Thing, error)
}

type ThingService struct {
	Things *shared.ThingStore
	l      log.Logger
}

func NewThingSvc(l log.Logger) ThingSvc {
	return ThingService{shared.SeedThings(), l}
}

func (t ThingService) GetAThing(_ context.Context, id string) (*shared.Thing, error) {
	return t.Things.Find(id)
}

func (t ThingService) GetAllThings(_ context.Context) ([]shared.Thing, error) {
	return t.Things.GetAllThings()
}
