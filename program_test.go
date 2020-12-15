package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	open := func(name string) (io.ReadCloser, error) {
		return ioutil.NopCloser(strings.NewReader(`[{
			"text": "coding"}]`)), nil
	}
	p := &Program{
		open: open,
		dir:  "/",
	}
	p.Load()
	if actual := p.records[0].Text; actual != "coding" {
		t.Errorf("expected text coding got %s", actual)
	}
}

func TestLoadError(t *testing.T) {
	someError := errors.New("some error")
	p := &Program{dir: "/"}
	p.open = func(string) (io.ReadCloser, error) { return nil, someError }
	if actual := p.Load(); actual != someError {
		t.Errorf("expected error to be %v got %v", someError, actual)
	}
}

type nopCloser struct {
	io.ReadWriter
}

func (nopCloser) Close() error { return nil }

func TestSave(t *testing.T) {
	b := nopCloser{&bytes.Buffer{}}
	create := func(name string) (io.WriteCloser, error) {
		return b, nil
	}
	p := &Program{
		create: create,
		dir:    "/",
	}
	p.Add(time.Now(), "foo")
	p.Save()
	dec := json.NewDecoder(b)
	records := make([]Record, 0)
	if err := dec.Decode(&records); err != nil {
		log.Fatal(err)
	}
	expect := p.records[0]
	actual := records[0]
	if !actual.Start.Equal(expect.Start) {
		t.Errorf("expected start to be %s got %s", expect.Start, actual.Start)
	}
	if !actual.Stop.Equal(expect.Stop) {
		t.Errorf("expect stop to be %s got %s", expect.Stop, actual.Stop)
	}
	if actual.Text != expect.Text {
		t.Errorf("expect text to be %s got %s", expect.Text, actual.Text)
	}
}

func TestSaveError(t *testing.T) {
	someError := errors.New("some error")
	p := &Program{dir: "/"}
	p.create = func(string) (io.WriteCloser, error) { return nil, someError }
	if actual := p.Save(); actual != someError {
		t.Errorf("expected error to be %v got %v", someError, actual)
	}
}

func TestName(t *testing.T) {
	td := []struct {
		name   string
		dir    string
		env    string
		usr    string
		expect string
	}{
		{"dir", "/dir", "", "", "/dir/20201221"},
		{"env", "", "/env", "", "/env/20201221"},
		{"usr", "", "", "/usr", "/usr/.td/20201221"},
	}
	for _, d := range td {
		t.Run(d.name, func(t *testing.T) {
			p := &Program{
				dir:     d.dir,
				envDir:  d.env,
				usrDir:  func() string { return d.usr },
				started: time.Date(2020, 12, 21, 23, 1, 0, 0, time.Local),
			}
			actual := p.Name()
			if actual != d.expect {
				t.Errorf("expect name to be %s got %s", d.expect, actual)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	p := &Program{dir: "/"}
	m := time.Now()
	expect := []Record{
		{Start: m, Text: "coding", Stop: m.Add(time.Hour)},
		{Start: m.Add(time.Hour), Text: "reading"},
	}
	for _, r := range expect {
		p.Add(r.Start, r.Text)
	}
	for i, r := range expect {
		t.Run(r.Text, func(t *testing.T) {
			actual := p.records[i]
			if !actual.Start.Equal(r.Start) {
				t.Errorf("expected start to be %v got %v", r.Start, actual.Start)
			}
			if !actual.Stop.Equal(r.Stop) {
				t.Errorf("expected stop to be %v got %v", r.Stop, actual.Stop)
			}
			if actual.Text != r.Text {
				t.Errorf("expected text to be %v got %v", r.Text, actual.Text)
			}
		})
	}
}

func TestNewStart(t *testing.T) {
	m := time.Now()
	td := []struct {
		name   string
		m      time.Time
		d      time.Duration
		expect time.Time
	}{
		{"negative duration", m, -1, m},
		{"zero duration", m, 0, time.Date(m.Year(), m.Month(), m.Day(), 0, 0, 0, 0, m.Location())},
		{"1h duration", m, time.Hour, time.Date(m.Year(), m.Month(), m.Day(), 1, 0, 0, 0, m.Location())},
		{"23h59m duration", m, 23*time.Hour + 59*time.Minute, time.Date(m.Year(), m.Month(), m.Day(), 23, 59, 0, 0, m.Location())},
		{"24h duration", m, 24 * time.Hour, m},
	}
	for _, d := range td {
		t.Run(d.name, func(t *testing.T) {
			actual := newStart(d.m, d.d)
			if !actual.Equal(d.expect) {
				t.Errorf("expected start to be %v got %v", d.expect, actual)
			}
		})
	}
}
