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
	create func(name string) (io.WriteCloser, error)
	dir    string
	envDir string
	name   string
	open   func(name string) (io.ReadCloser, error)
	usrDir func() string

	started time.Time

	exactly          string
	prefix           string
	records          []Record
	roundDuration    time.Duration
	truncateDuration time.Duration
}

func (p *Program) Load() {
	f, err := p.open(p.Name())
	if err != nil && !os.IsNotExist(err) {
		logger.Println(err)
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&p.records)
	if err != nil {
		logger.Println(err)
	}
}

func (p *Program) Save() {
	if _, err := os.Stat(p.Dir()); os.IsNotExist(err) {
		_ = os.Mkdir(p.dir, 0700)
	}
	w, err := p.create(p.Name())
	if err != nil {
		logger.Fatalln(err)
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(p.records)
	if err != nil {
		logger.Fatalln(err)
	}
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

func (p *Program) Print(t time.Time, w io.Writer) {
	fmt.Println(p.started.Format("Monday, January 2, 2006"))
	td := time.Duration(0)
	for _, r := range p.records {
		d := r.Dur(t)
		pre := " "
		b := p.match(r)
		if b {
			pre = "+"
			td += d
		}
		fmt.Print(pre)
		fmt.Printf(
			"%v %v %v\n",
			p.formatRange(r),
			p.formatDur(d),
			r.Text)
	}
	fmt.Println(p.formatDur(td))
}

func (p *Program) match(r Record) bool {
	return p.exactly != "" && p.exactly == r.Text ||
		p.prefix != "" && strings.HasPrefix(r.Text, p.prefix)
}

func (p *Program) formatRange(r Record) string {
	format := func(t time.Time) string {
		const layout = "15:04"
		if t.IsZero() {
			return strings.Repeat(" ", len(layout))
		}
		return t.Format(layout)
	}
	return fmt.Sprintf("%v – %v", format(r.Start), format(r.Stop))
}

func (p *Program) formatDur(d time.Duration) string {
	d = d.Round(p.roundDuration)
	d = d.Truncate(p.truncateDuration)
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}
