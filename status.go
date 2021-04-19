package main

type Status int

const (
	Initial Status = iota
	OK
	Alert
	Unknown
)

func (s *Status) String() string {
	switch {
	case *s == Initial:
		return "Initial"
	case *s == OK:
		return "OK"
	case *s == Alert:
		return "Alert"
	default:
		return "Unknown"
	}
}

func (s *Status) Recovery() {
	*s = OK
}

func (s *Status) Trigger() {
	*s = Alert
}

func (s *Status) Unknown() {
	*s = Unknown
}

func (s *Status) Is(status Status) bool {
	return *s == status
}
