package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var (
	logger  = log.New(os.Stdout, "td: ", log.Lshortfile)
	pattern = flag.String("pattern", "", "")
	start   = flag.Duration("start", time.Duration(-1), "")
)

func main() {
	p := &Program{
		create:  create,
		envDir:  os.Getenv("TDPATH"),
		open:    open,
		pre:     "\033[1m",
		started: time.Now(),
		suf:     "\033[0m",
		usrDir:  usrDir,
	}
	flag.BoolVar(&p.printDuration, "printDuration", true, "")
	flag.BoolVar(&p.printRange, "printRange", false, "")
	flag.DurationVar(&p.offset, "offset", time.Duration(0), "")
	flag.DurationVar(&p.roundDur, "roundDuration", time.Duration(time.Minute), "")
	flag.DurationVar(&p.truncateDur, "truncateDuration", time.Duration(0), "")
	flag.StringVar(&p.dir, "homeDir", "", "home directory, if not set $TDPATH or $HOME/.td is used")
	flag.Parse()

	if *pattern != "" {
		p.pattern = regexp.MustCompile(*pattern)
	}
	if err := p.Load(); err != nil {
		logger.Fatal(err)
	}
	t := newStart(p.started, *start)
	if p.offset != 0 {
		t = t.Add(p.offset)
	}
	if flag.NArg() == 0 {
		p.Print(t)
		return
	}
	cmd, text := cmdText(flag.Args())
	switch cmd {
	case "start":
		p.Add(t, text)
		if err := p.Save(); err != nil {
			logger.Fatal(err)
		}
	}
}

func open(name string) (r io.ReadCloser, err error) {
	r, err = os.Open(name)
	return
}

func create(name string) (w io.WriteCloser, err error) {
	dir, _ := path.Split(name)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)
	}
	if err == nil {
		w, err = os.Create(name)
	}
	return
}

func usrDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatalln(err)
	}
	return dir
}

func newStart(t time.Time, d time.Duration) time.Time {
	if d < 0 || d >= 24*time.Hour {
		return t
	}
	return newTime(t, d)
}

func newTime(t time.Time, d time.Duration) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0,
		0,
		0,
		0,
		t.Location()).
		Add(d)
}

func cmdText(args []string) (cmd, text string) {
	if len(args) == 0 {
		return
	}
	cmd = args[0]
	text = strings.Join(args[1:], " ")
	return
}
