package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_String(t *testing.T) {
	cases := []struct {
		status Status
		want   string
	}{
		{
			status: OK,
			want:   "OK"},
		{
			status: Critical,
			want:   "CRITICAL",
		},
		{
			status: Unknown,
			want:   "UNKNOWN",
		},
	}

	for _, c := range cases {
		t.Run(c.status.String(), func(t *testing.T) {
			got := c.status.String()
			assert.Equal(t, c.want, got)
		})
	}
}

func TestStatus_Recovery(t *testing.T) {
	for _, s := range []Status{OK, Critical, Unknown} {
		s.Recovery()
		assert.Equal(t, OK, s)
	}
}

func TestStatus_Trigger(t *testing.T) {
	for _, s := range []Status{OK, Critical, Unknown} {
		s.Trigger()
		assert.Equal(t, Critical, s)
	}
}

func TestStatus_Unknown(t *testing.T) {
	for _, s := range []Status{OK, Critical, Unknown} {
		s.Unknown()
		assert.Equal(t, Unknown, s)
	}
}

func TestStatus_Is_true(t *testing.T) {
	cases := []struct {
		status Status
		is     Status
	}{
		{
			status: OK,
			is:     OK,
		},
		{
			status: Critical,
			is:     Critical,
		},
		{
			status: Unknown,
			is:     Unknown,
		},
	}

	for _, c := range cases {
		t.Run(c.status.String(), func(t *testing.T) {
			assert.True(t, c.status.Is(c.is))
		})
	}
}

func TestStatus_Is_false(t *testing.T) {
	cases := []struct {
		status Status
		is     Status
	}{
		{
			status: OK,
			is:     Critical,
		},
		{
			status: OK,
			is:     Unknown,
		},
		{
			status: Critical,
			is:     OK,
		},
		{
			status: Critical,
			is:     Unknown,
		},
		{
			status: Unknown,
			is:     OK,
		},
		{
			status: Unknown,
			is:     Critical,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s->%s", c.status.String(), c.is.String()), func(t *testing.T) {
			assert.False(t, c.status.Is(c.is))
		})
	}
}
