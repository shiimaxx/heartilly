package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlackNotifier_color(t *testing.T) {
	cases := []struct {
		status Status
		want   string
	}{
		{status: OK, want: "good"},
		{status: Critical, want: "danger"},
		{status: Unknown, want: "#808080"},
	}

	s := SlackNotifier{}
	for _, c := range cases {
		t.Run(c.status.String(), func(t *testing.T) {
			got := s.color(c.status)
			assert.Equal(t, c.want, got)
		})
	}
}
