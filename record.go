package main

import (
	"time"
)

type Record struct {
	Start time.Time
	Stop  time.Time
	Text  string
}

func (r Record) Dur(t time.Time) time.Duration {
	if r.Stop.IsZero() {
		return t.Sub(r.Start)
	}
	return r.Stop.Sub(r.Start)
}
