package main

import "time"

func Example() {
	t := time.Date(2020, 12, 21, 19, 30, 0, 0, time.Local)
	p := &Program{
		pre:           ">",
		prefix:        "foo",
		printDuration: true,
		printRange:    true,
		started:       t,
		suf:           "<",
	}
	p.Add(t, "foo")
	p.Add(t.Add(15*time.Minute), "bar")
	p.Add(t.Add(30*time.Minute), "foo")
	p.Print(t.Add(time.Hour))
	// Output:
	// Monday, December 21, 2020
	// >19:30 – 19:45    15m foo<
	// 19:45 – 20:00    15m bar
	// >20:00 –          30m foo<
	// --------------------
	// >                 45m<
}
