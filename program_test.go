package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
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
	td := []struct {
		name        string
		returnError error
		expectError error
	}{
		{"return error", someError, someError},
		{"return nil for ErrNotExist", os.ErrNotExist, nil},
	}
	p := &Program{dir: "/"}
	for _, d := range td {
		t.Run(d.name, func(t *testing.T) {
			p.open = func(string) (io.ReadCloser, error) { return nil, d.returnError }
			if actual := p.Load(); actual != d.expectError {
				t.Errorf("expected error to be %v got %v", d.expectError, actual)
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
