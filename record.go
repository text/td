package main

import "time"

type Record struct {
	Started  time.Time
	duration time.Duration
}

func (r Record) IsStarted() bool { return !r.Started.IsZero() }

func (r Record) Duration(t time.Time) (d time.Duration) {
	d += r.duration
	if r.IsStarted() {
		d += t.Sub(r.Started)
	}
	return
}
