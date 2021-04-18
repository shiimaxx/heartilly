package main

type Status int

const (
	OK Status = iota
	ALERT
)

func (s *Status) String() string {
	switch {
	case *s == OK:
		return "OK"
	case *s == ALERT:
		return "Alert"
	default:
		return "Unknown"
	}
}

func (s *Status) Recovery() {
	*s = OK
}

func (s *Status) Trigger() {
	*s = ALERT
}
