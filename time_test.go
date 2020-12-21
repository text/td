package main

import "time"

func fakeTime(hour, min int) time.Time {
	return time.Date(
		2020,
		12,
		21,
		hour,
		min,
		0,
		0,
		time.Local)
}
