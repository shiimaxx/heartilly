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
			status: OK,
			want: "OK", },
		{
			status: Critical,
			want: "CRITICAL",
		},
		{
			status: Unknown,
			want: "UNKNOWN",
		},
	}

	for _, c := range cases {
		got := c.status.String()
		assert.Equal(t, c.want, got)
	}
}

func TestRecovery(t *testing.T) {
	for _, s := range []Status{OK, Critical, Unknown} {
		s.Recovery()
		assert.Equal(t, OK, s)
	}
}

func TestTrigger(t *testing.T) {
	for _, s := range []Status{OK, Critical, Unknown} {
		s.Trigger()
		assert.Equal(t, Critical, s)
	}
}

func TestUnknown(t *testing.T) {
	for _, s := range []Status{OK, Critical, Unknown} {
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
			status: OK,
			is: OK,
		},
		{
			status: Critical,
			is: Critical,
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
			status: OK,
			is: Critical,
		},
		{
			status: OK,
			is: Unknown,
		},
		{
			status: Critical,
			is: OK,
		},
		{
			status: Critical,
			is: Unknown,
		},
		{
			status: Unknown,
			is: OK,
		},
		{
			status: Unknown,
			is: Critical,
		},
	}

	for _, c := range cases {
		assert.False(t, c.status.Is(c.is))	
	}
}
