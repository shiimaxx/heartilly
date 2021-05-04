package main

import "time"

type Result struct {
	ID        int64     `json:"id" db:"id"`
	CheckedAt time.Time `json:"checked_at" db:"checked_at"`
	Status    string    `json:"status" db:"status"`
	Reason    string    `json:"reason" db:"reason"`
	MonitorID int64     `json:"-" db:"monitor_id"`
}
