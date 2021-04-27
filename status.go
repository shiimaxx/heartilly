package main

type Status int

const (
	OK Status = iota
	Critical
	Unknown
)

func (s *Status) String() string {
	switch {
	case *s == OK:
		return "OK"
	case *s == Critical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func (s *Status) Recovery() {
	*s = OK
}

func (s *Status) Trigger() {
	*s = Critical
}

func (s *Status) Unknown() {
	*s = Unknown
}

func (s *Status) Is(status Status) bool {
	return *s == status
}
