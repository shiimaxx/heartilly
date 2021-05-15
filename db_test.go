package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
)

func prepareTestDB(t *testing.T) func() {
	t.Helper()

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create tempdir failed:", err)
	}

	dbfile := fmt.Sprintf("%s/heartilly_test.db", dir)
	err = OpenDB(dbfile)
	if err != nil {
		t.Fatal("open db failed:", err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("sqlite"),
		testfixtures.Directory("testdata/fixtures"),
	)
	if err != nil {
		t.Fatal("create test fixtures failed:", err)
	}

	if err := fixtures.Load(); err != nil {
		t.Fatal("load fixtures failed:", err)
	}

	return func() { os.RemoveAll(dir) }
}

func TestOpenDB(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create tempdir failed:", err)
	}
	defer os.RemoveAll(dir)

	dbfile := fmt.Sprintf("%s/heartilly_test.db", dir)
	err = OpenDB(dbfile)

	assert.Nil(t, err)

	_, err = db.Queryx(`SELECT * FROM sqlite_master WHERE name = "monitor"`)
	assert.Nil(t, err)

	_, err = db.Queryx(`SELECT * FROM sqlite_master WHERE name = "result"`)
	assert.Nil(t, err)
}

func TestGetMonitors(t *testing.T) {
	cleanup := prepareTestDB(t)
	defer cleanup()

	got, err := GetAllMonitors()

	assert.Nil(t, err)

	want := []*Monitor{
		{
			ID:     1,
			Name:   "GET /monitor/get",
			Method: "GET",
			URL:    parseURL(t, "http://example.com/monitor/get"),
			Follow: false,
		},
		{
			ID:     2,
			Name:   "POST /monitor/post",
			Method: "POST",
			URL:    parseURL(t, "http://example.com/monitor/post"),
			Follow: false,
		},
		{
			ID:     3,
			Name:   "GET /monitor/follow",
			Method: "GET",
			URL:    parseURL(t, "http://example.com/monitor/follow"),
			Follow: true,
		},
	}

	assert.Equal(t, want, got)
}

func TestGetMonitorByName(t *testing.T) {
	cleanup := prepareTestDB(t)
	defer cleanup()

	cases := []struct {
		name string
		want *Monitor
	}{
		{
			name: "GET /monitor/get",
			want: &Monitor{
				ID:     1,
				Name:   "GET /monitor/get",
				Method: "GET",
				URL:    parseURL(t, "http://example.com/monitor/get"),
				Follow: false,
			},
		},
		{
			name: "POST /monitor/post",
			want: &Monitor{
				ID:     2,
				Name:   "POST /monitor/post",
				Method: "POST",
				URL:    parseURL(t, "http://example.com/monitor/post"),
				Follow: false,
			},
		},
		{
			name: "GET /monitor/follow",
			want: &Monitor{
				ID:     3,
				Name:   "GET /monitor/follow",
				Method: "GET",
				URL:    parseURL(t, "http://example.com/monitor/follow"),
				Follow: true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := GetMonitorByName(c.name)

			assert.Nil(t, err)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestGetResults(t *testing.T) {
	cleanup := prepareTestDB(t)
	defer cleanup()

	baseTime, err := time.Parse("2006-01-02 15:04:05 +0000", "2006-01-02 15:04:05 +0000")
	if err != nil {
		t.Fatal("parse time failed:", err)
	}

	cases := []struct {
		name      string
		monitorID int
		want      []*Result
	}{
		{
			name:      "monitorID: 1",
			monitorID: 1,
			want: []*Result{
				{
					ID:        1,
					CheckedAt: baseTime,
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 1,
				},
				{
					ID:        4,
					CheckedAt: baseTime.Add(1 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 1,
				},
				{
					ID:        7,
					CheckedAt: baseTime.Add(2 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 1,
				},
				{
					ID:        10,
					CheckedAt: baseTime.Add(3 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 1,
				},
				{
					ID:        13,
					CheckedAt: baseTime.Add(4 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 1,
				},
			},
		},
		{
			name:      "monitorID: 2",
			monitorID: 2,
			want: []*Result{
				{
					ID:        2,
					CheckedAt: baseTime.Add(10 * time.Second),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 2,
				},
				{
					ID:        5,
					CheckedAt: baseTime.Add(10 * time.Second).Add(1 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 2,
				},
				{
					ID:        8,
					CheckedAt: baseTime.Add(10 * time.Second).Add(2 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 2,
				},
				{
					ID:        11,
					CheckedAt: baseTime.Add(10 * time.Second).Add(3 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 2,
				},
				{
					ID:        14,
					CheckedAt: baseTime.Add(10 * time.Second).Add(4 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 2,
				},
			},
		},
		{
			name:      "monitorID: 3",
			monitorID: 3,
			want: []*Result{
				{
					ID:        3,
					CheckedAt: baseTime.Add(20 * time.Second),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 3,
				},
				{
					ID:        6,
					CheckedAt: baseTime.Add(20 * time.Second).Add(1 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 3,
				},
				{
					ID:        9,
					CheckedAt: baseTime.Add(20 * time.Second).Add(2 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 3,
				},
				{
					ID:        12,
					CheckedAt: baseTime.Add(20 * time.Second).Add(3 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 3,
				},
				{
					ID:        15,
					CheckedAt: baseTime.Add(20 * time.Second).Add(4 * time.Minute),
					Status:    "OK",
					Reason:    "200 OK",
					MonitorID: 3,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := GetResultsByMonitorID(c.monitorID)

			assert.Nil(t, err)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestCreateResult(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create tempdir failed:", err)
	}
	defer os.RemoveAll(dir)

	dbfile := fmt.Sprintf("%s/heartilly_test.db", dir)
	if err = OpenDB(dbfile); err != nil {
		t.Fatal("open db failed:", err)
	}

	checkedAt, err := time.Parse("2006-01-02 15:04:05 +0000", "2006-01-02 15:04:05 +0000")
	if err != nil {
		t.Fatal("parse time failed:", err)
	}

	result := &Result{
		CheckedAt: checkedAt,
		Status:    "OK",
		Reason:    "200 OK",
		MonitorID: 1,
	}
	err = CreateResult(result)

	assert.Nil(t, err)

	want := Result{
		ID:        1,
		CheckedAt: checkedAt,
		Status:    "OK",
		Reason:    "200 OK",
		MonitorID: 1,
	}
	got := Result{}

	if err := db.Get(&got, `SELECT * FROM result`); err != nil {
		t.Fatal("query failed:", err)
	}

	assert.Equal(t, want, got)

}
