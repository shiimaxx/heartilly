package main

import (
	"time"
)

type Result struct {
	ID      int64     `json:"id" db:"id"`
	Created time.Time `json:"created" db:"created"`
	Status  string    `json:"status" db:"status"`
	Reason  string    `json:"reason" db:"reason"`
}
