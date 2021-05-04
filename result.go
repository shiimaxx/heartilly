package main

type Result struct {
	ID        int64  `json:"id" db:"id"`
	Created   uint64 `json:"created" db:"created"`
	Status    string `json:"status" db:"status"`
	Reason    string `json:"reason" db:"reason"`
	MonitorID int64  `json:"-" db:"monitor_id"`
}
