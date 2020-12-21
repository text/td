package main

import (
	"io"
)

type option func(*Program) option

func (p *Program) Option(opts ...option) (previous option) {
	for _, opt := range opts {
		previous = opt(p)
	}
	return previous
}

// OutputWriter sets Program's output writer to w.
func OutputWriter(w io.Writer) option {
	return func(p *Program) option {
		previous := p.outputWriter
		p.outputWriter = w
		return OutputWriter(previous)
	}
}
