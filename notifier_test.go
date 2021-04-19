package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor(t *testing.T) {
	cases := []struct {
		status Status
		want   string
	}{
		{status: OK, want: "good"},
		{status: Alert, want: "danger"},
		{status: Unknown, want: "#808080"},
	}

	s := SlackNotifier{}
	for _, c := range cases {
		got := s.color(c.status)
		assert.Equal(t, c.want, got)
	}
}
