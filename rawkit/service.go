package rawkit

import (
	"context"
	"using-kit/shared"
)

type StrSvc interface {
	GetAThing(context.Context, string) (*shared.Thing, error)
	GetAllThings(context.Context) ([]shared.Thing, error)
}

type ThingService struct {
	Things *shared.ThingStore
}

func (t ThingService) GetAThing(_ context.Context, id string) (*shared.Thing, error) {
	return t.Things.Find(id)
}

func (t ThingService) GetAllThings(c context.Context) ([]shared.Thing, error) {
	return t.Things.GetAllThings()
}

//In Go kit, the primary messaging pattern is RPC.
//So, every method in our interface will be modeled as a remote procedure call.
//For each method, we define request and response structs,
//capturing all of the input and output parameters respectively.
type getThingRequest struct {
	ID string `json:"id"`
}

type getThingResponse struct {
	Thing shared.Thing `json:"thing"`
	Err   string       `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type allThingsResponse struct {
	Ts []shared.Thing `json:"v"`
}
