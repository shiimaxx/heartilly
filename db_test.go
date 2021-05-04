package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
)

func TestOpenDB(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create tempdir failed: ", err)
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
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create tempdir failed: ", err)
	}
	defer os.RemoveAll(dir)

	dbfile := fmt.Sprintf("%s/heartilly_test.db", dir)
	err = OpenDB(dbfile)
	if err != nil {
		t.Fatal("open db failed: ", err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("sqlite"),
		testfixtures.Directory("testdata/fixtures"),
	)
	if err != nil {
		t.Fatal("create test fixtures failed: ", err)
	}

	if err := fixtures.Load(); err != nil {
		t.Fatal("load fixtures failed: ", err)
	}

	got, err := GetMonitors()
	if err != nil {
		t.Fatal("get monitors failed: ", err)
	}

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
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create tempdir failed: ", err)
	}
	defer os.RemoveAll(dir)

	dbfile := fmt.Sprintf("%s/heartilly_test.db", dir)
	err = OpenDB(dbfile)
	if err != nil {
		t.Fatal("open db failed: ", err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("sqlite"),
		testfixtures.Directory("testdata/fixtures"),
	)
	if err != nil {
		t.Fatal("create test fixtures failed: ", err)
	}

	if err := fixtures.Load(); err != nil {
		t.Fatal("load fixtures failed: ", err)
	}

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
			if err != nil {
				t.Fatal("get monitors failed: ", err)
			}

			assert.Equal(t, c.want, got)
		})
	}
}

func TestGetResults(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create tempdir failed: ", err)
	}
	defer os.RemoveAll(dir)

	dbfile := fmt.Sprintf("%s/heartilly_test.db", dir)
	err = OpenDB(dbfile)
	if err != nil {
		t.Fatal("open db failed: ", err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("sqlite"),
		testfixtures.Directory("testdata/fixtures"),
	)
	if err != nil {
		t.Fatal("create test fixtures failed: ", err)
	}

	if err := fixtures.Load(); err != nil {
		t.Fatal("load fixtures failed: ", err)
	}

	got, err := GetResults(1)
	if err != nil {
		t.Fatal("get results failed: ", err)
	}

	want := []*Result{
		{ID: 1, Created: 1136239445, Status: "OK", Reason: "200 OK", MonitorID: 1},
		{ID: 4, Created: 1136239505, Status: "OK", Reason: "200 OK", MonitorID: 1},
		{ID: 7, Created: 1136239565, Status: "OK", Reason: "200 OK", MonitorID: 1},
		{ID: 10, Created: 1136239625, Status: "OK", Reason: "200 OK", MonitorID: 1},
		{ID: 13, Created: 1136239685, Status: "OK", Reason: "200 OK", MonitorID: 1},
	}

	assert.Equal(t, want, got)
}
