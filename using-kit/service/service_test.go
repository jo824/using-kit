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
