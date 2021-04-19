package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	cases := []struct {
		status Status
		want   string
	}{
		{
			status: Initial,
			want: "Initial",
		},
		{
			status: OK,
			want: "OK", },
		{
			status: Alert,
			want: "Alert",
		},
		{
			status: Unknown,
			want: "Unknown",
		},
	}

	for _, c := range cases {
		got := c.status.String()
		assert.Equal(t, c.want, got)
	}
}

func TestRecovery(t *testing.T) {
	for _, s := range []Status{Initial, OK, Alert, Unknown} {
		s.Recovery()
		assert.Equal(t, OK, s)
	}
}

func TestTrigger(t *testing.T) {
	for _, s := range []Status{Initial, OK, Alert, Unknown} {
		s.Trigger()
		assert.Equal(t, Alert, s)
	}
}

func TestUnknown(t *testing.T) {
	for _, s := range []Status{Initial, OK, Alert, Unknown} {
		s.Unknown()
		assert.Equal(t, Unknown, s)
	}
}

func TestIs_true(t *testing.T) {
	cases := []struct {
		status Status
		is Status
	}{
		{
			status: Initial,
			is: Initial,
		},
		{
			status: OK,
			is: OK,
		},
		{
			status: Alert,
			is: Alert,
		},
		{
			status: Unknown,
			is: Unknown,
		},
	}

	for _, c := range cases {
		assert.True(t, c.status.Is(c.is))	
	}
}

func TestIs_false(t *testing.T) {
	cases := []struct {
		status Status
		is Status
	}{
		{
			status: Initial,
			is: OK,
		},
		{
			status: Initial,
			is: Alert,
		},
		{
			status: Initial,
			is: Unknown,
		},
		{
			status: OK,
			is: Initial,
		},
		{
			status: OK,
			is: Alert,
		},
		{
			status: OK,
			is: Unknown,
		},
		{
			status: Alert,
			is: Initial,
		},
		{
			status: Alert,
			is: OK,
		},
		{
			status: Alert,
			is: Unknown,
		},
		{
			status: Unknown,
			is: Initial,
		},
		{
			status: Unknown,
			is: OK,
		},
		{
			status: Unknown,
			is: Alert,
		},
	}

	for _, c := range cases {
		assert.False(t, c.status.Is(c.is))	
	}
}
