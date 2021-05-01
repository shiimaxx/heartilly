package main

import (
	"database/sql"
	"database/sql/driver"
	"net/url"
)

type Monitor struct {
	ID     int64  `json:"id" toml:"-" db:"id"`
	Name   string `json:"name" toml:"name" db:"name"`
	Method string `json:"method" toml:"method" db:"method"`
	URL    URL    `json:"url" toml:"url" db:"url"`
	Follow bool   `json:"follow" toml:"follow" db:"follow"`
}

func InitSyncMonitor(monitors []*Monitor) ([]*Monitor, error) {
	var notFound []*Monitor

	for _, m := range monitors {
		_, err := GetMonitorByName(m.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				notFound = append(notFound, m)
			} else {
				return nil, err
			}
		}
	}

	if err := CreateMonitors(notFound); err != nil {
		return nil, err
	}

	return GetMonitors()
}

type URL url.URL

// https://golang.org/pkg/database/sql/driver/#Value
func (u *URL) Value() (driver.Value, error) {
	return driver.Value(u.String()), nil
}

// https://golang.org/pkg/database/sql/#Scanner
func (u *URL) Scan(value interface{}) error {
	parsedURL, err := url.Parse(value.(string))
	if err != nil {
		return err
	}

	*u = URL(*parsedURL)

	return nil
}

// https://golang.org/pkg/encoding/#TextMarshaler
func (u *URL) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// https://golang.org/pkg/encoding/#TextUnmarshaler
func (u *URL) UnmarshalText(text []byte) error {
	parsedURL, err := url.Parse(string(text))
	if err != nil {
		return err
	}

	*u = URL(*parsedURL)

	return nil
}

func (u *URL) String() string {
	return (*url.URL)(u).String()
}
