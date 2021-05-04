package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func OpenDB(dbfile string) error {
	var err error
	db, err = sqlx.Open("sqlite3", dbfile)
	if err != nil {
		return err
	}

	createMonitor := `
	CREATE TABLE IF NOT EXISTS monitor (
	  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  name TEXT UNIQUE,
	  method TEXT,
	  url TEXT, 
	  follow INTEGER
	);
	`
	if _, err := db.Exec(createMonitor); err != nil {
		return err
	}

	createResult := `
	CREATE TABLE IF NOT EXISTS result (
	  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  checked_at TIMESTAMP,
	  status TEXT,
	  reason TEXT,
	  monitor_id INTEGER,
	  FOREIGN KEY(monitor_id) REFERENCES monitor(id)
	);
	`
	if _, err := db.Exec(createResult); err != nil {
		return err
	}
	return nil
}

func GetMonitors() ([]*Monitor, error) {
	var monitors []*Monitor

	query := `SELECT * FROM monitor`
	if err := db.Select(&monitors, query); err != nil {
		return nil, err
	}

	return monitors, nil
}

func GetMonitorByName(name string) (*Monitor, error) {
	query := `SELECT * FROM monitor WHERE name = ?`
	monitor := Monitor{}

	if err := db.Get(&monitor, query, name); err != nil {
		return nil, err
	}

	return &monitor, nil
}

func CreateMonitors(monitors []*Monitor) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO monitor(name, method, url, follow) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	for _, m := range monitors {
		_, err = stmt.Exec(m.Name, m.Method, m.URL.String(), m.Follow)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetResults(id int) ([]*Result, error) {
	var results []*Result
	query := `SELECT * FROM result WHERE monitor_id = ?`

	if err := db.Select(&results, query, id); err != nil {
		return nil, err
	}

	return results, nil
}

func CreateResult(result *Result) error {
	query := `INSERT INTO result(checked_at, status, reason, monitor_id) VALUES(?, ?, ?, ?)`

	_, err := db.Exec(query, result.CheckedAt, result.Status, result.Reason, result.MonitorID)
	if err != nil {
		return err
	}

	return nil
}
