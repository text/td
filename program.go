package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Program struct {
	prevName     string
	data         map[string]Record
	outputWriter io.Writer
}

func NewProgram() *Program {
	return &Program{
		data: make(map[string]Record),
	}
}

func (p *Program) ProcessArguments(t time.Time, args []string) {
	if len(args) == 0 {
		p.output(t)
		return
	}
	command := args[0]
	switch command {
	case "start":
		p.start(t, strings.Join(args[1:], " "))
	}
}

func (p *Program) start(t time.Time, name string) {
	if prev, ok := p.data[p.prevName]; ok {
		prev.duration += t.Sub(prev.Started)
		prev.Started = time.Time{}
		p.data[p.prevName] = prev
	}
	p.data[name] = Record{
		Started: t,
	}
	p.prevName = name
}

func (p *Program) output(t time.Time) error {
	fmt.Fprintln(
		p.outputWriter,
		t.Format("Monday, January 2, 2006"))
	for name, value := range p.data {
		d := value.Duration(t)
		fmt.Fprintf(
			p.outputWriter,
			"%9v %s\n",
			formatDuration(d),
			name)
	}
	return nil
}

func formatDuration(d time.Duration) string {
	s := d.String()
	s = strings.Replace(s, "h0m0s", "h      ", 1)
	s = strings.Replace(s, "m0s", "m   ", 1)
	return s
}
