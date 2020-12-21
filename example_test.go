package main

import "os"

func Example() {
	p := NewProgram()
	p.Option(OutputWriter(os.Stdout))
	p.ProcessArguments(fakeTime(16, 0), []string{"start", "foo"})
	p.ProcessArguments(fakeTime(17, 0), []string{"start", "bar"})
	p.ProcessArguments(fakeTime(17, 15), nil)
	// Output:
	// Monday, December 21, 2020
	//  1h       foo
	//    15m    bar
}
