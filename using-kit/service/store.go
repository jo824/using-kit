package service

import "errors"

type Thing struct {
	ID        string
	Available bool
}

type ThingStore struct {
	things map[string]*Thing
}

func SeedThings() *ThingStore {
	tings := map[string]*Thing{
		"abd": {
			ID:        "abd",
			Available: true,
		},
		"eek": {
			ID:        "eek",
			Available: true,
		},
		"yik": {
			ID:        "yik",
			Available: true,
		},
		"yak": {
			ID:        "yak",
			Available: true,
		},
	}
	return &ThingStore{tings}
}

func (ts *ThingStore) Save(t *Thing) error {
	if _, ok := ts.things[t.ID]; ok {
		return ErrAlreadyExists
	}
	ts.things[t.ID] = t
	return nil
}

func (ts *ThingStore) Find(id string) (*Thing, error) {
	if val, ok := ts.things[id]; ok {
		return val, nil
	}
	return nil, errors.New("id not found")
}

func (ts *ThingStore) GetAllThings() ([]Thing, error) {
	tls := make([]Thing, 0, len(ts.things))
	for _, v := range ts.things {
		tls = append(tls, *v)
	}
	return tls, nil
}
