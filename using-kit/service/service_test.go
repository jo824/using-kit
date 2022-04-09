package service

import (
	"context"
	"testing"
)

func TestGetAThing(t *testing.T) {
	tt := []struct {
		name     string
		id       string
		expected bool
	}{
		{
			"first",
			"yik",
			true,
		},
		{
			"second",
			"yak",
			true,
		},
		{"third",
			"nope",
			false,
		},
	}
	s := NewThingSvc()
	for _, tc := range tt {
		t.Run(tc.id, func(t *testing.T) {
			res, _ := s.GetAThing(context.TODO(), tc.id)
			if res == nil && tc.expected {
				t.Fail()
				t.Logf("Value should exist for test %s", tc.name)
			}
		})
	}
}

func TestGetAllThings(t *testing.T) {
	expected := 4

	svc := NewThingSvc()
	result, err := svc.GetAllThings(context.TODO())
	if err != nil {
		t.Fail()
		t.Logf("GetAllThings failed with error: %s", err)
	}
	actual := len(result)
	if actual != expected {
		t.Fail()
		t.Logf("GetAllThings expected: %d, actual: %d", expected, actual)
	}

}

func TestAddThing(t *testing.T) {
	svc := NewThingSvc()

	th := &Thing{
		ID:        "go",
		Available: true,
	}

	err := svc.AddThing(context.TODO(), th)
	if err != nil {
		t.Fail()
		t.Logf("AddThing Failed error: %s", err)

	}
	err = svc.AddThing(context.TODO(), th)
	if err != ErrAlreadyExists {
		t.Fail()
		t.Logf("expect to fail with thing already exists error")
	}
}
