package main

import (
	"flag"
	"os"
	"time"
)

func main() {
	flag.Parse()
	t := time.Now()
	p := NewProgram()
	p.Option(OutputWriter(os.Stdout))
	p.ProcessArguments(t, flag.Args())
}
