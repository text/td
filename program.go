package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type Program struct {
	create        func(name string) (io.WriteCloser, error)
	dir           string
	envDir        string
	name          string
	open          func(name string) (io.ReadCloser, error)
	pre, suf      string
	prefix        string
	printDuration bool
	printRange    bool
	records       []Record
	roundDur      time.Duration
	started       time.Time
	truncateDur   time.Duration
	usrDir        func() string
}

func (p *Program) Load() error {
	f, err := p.open(p.Name())
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	return dec.Decode(&p.records)
}

func (p *Program) Save() error {
	w, err := p.create(p.Name())
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	return enc.Encode(p.records)
}

func (p *Program) Dir() string {
	if p.dir == "" {
		p.dir = p.envDir
	}
	if p.dir == "" {
		p.dir = path.Join(p.usrDir(), ".today")
	}
	return p.dir
}

func (p *Program) Name() string {
	if p.name == "" {
		name := p.started.Format("20060102")
		p.name = path.Join(p.Dir(), name)
	}
	return p.name
}

func (p *Program) Add(t time.Time, text string) {
	i := len(p.records) - 1
	if i >= 0 {
		r := p.records[i]
		r.Stop = t
		p.records[i] = r
	}
	p.records = append(p.records, Record{
		Start: t,
		Text:  text,
	})
}

func (p *Program) Print(t time.Time) {
	fmt.Println(p.started.Format("Monday, January 2, 2006"))
	td := time.Duration(0)
	for _, r := range p.records {
		d := r.Dur(t)
		include := p.include(r)
		if include {
			td += d
			fmt.Print(p.pre)
		}
		if p.printRange {
			fmt.Printf("%s ", p.formatRange(r))
		}
		if p.printDuration {
			fmt.Printf("%s ", p.formatDur(d))
		}
		fmt.Print(r.Text)
		if include {
			fmt.Print(p.suf)
		}
		fmt.Println()
	}
	if td > 0 {
		fmt.Println(strings.Repeat("-", p.seplen()))
		fmt.Print(p.pre)
		fmt.Print(strings.Repeat(" ", p.indent()))
		fmt.Print(p.formatDur(td))
		fmt.Print(p.suf)
		fmt.Println()
	}
}

func (p *Program) seplen() int {
	return p.indent() + len("21h59m")
}

func (p *Program) indent() int {
	if p.printRange {
		return len([]rune("15:04 – 15:04 "))
	}
	return 0
}

func (p *Program) include(r Record) bool {
	return p.prefix != "" && strings.HasPrefix(r.Text, p.prefix)
}

func (p *Program) formatRange(r Record) string {
	format := func(t time.Time) string {
		if t.IsZero() {
			return strings.Repeat(" ", len("15:04"))
		}
		return t.Format("15:04")
	}
	return fmt.Sprintf("%v – %v", format(r.Start), format(r.Stop))
}

func (p *Program) formatDur(d time.Duration) string {
	d = d.Round(p.roundDur)
	d = d.Truncate(p.truncateDur)
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return fmt.Sprintf("%6s", s)
}
