package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	logger = log.New(os.Stdout, "today: ", log.Lshortfile)
	p      *Program
)

func init() {
	p = &Program{
		create:  create,
		envDir:  os.Getenv("TODAYPATH"),
		open:    open,
		started: time.Now(),
		usrDir:  usrDir,
	}
	flag.DurationVar(&p.roundDuration, "roundDuration", time.Duration(0), "")
	flag.DurationVar(&p.truncateDuration, "truncateDuration", time.Duration(0), "")
	flag.StringVar(&p.dir, "homeDir", "", "home directory, if not set $TODAYPATH or $HOME/.today is used")
	flag.StringVar(&p.exactly, "exactly", "", "")
	flag.StringVar(&p.prefix, "hasPrefix", "", "")
}

func main() {
	flag.Parse()
	p.Load()
	text := strings.Join(flag.Args(), " ")
	if flag.NArg() == 0 {
		p.Print(p.started, os.Stdout)
		return
	}
	p.Add(time.Now(), text)
	p.Save()
}

func open(name string) (r io.ReadCloser, err error) {
	r, err = os.Open(name)
	return
}

func create(name string) (w io.WriteCloser, err error) {
	w, err = os.Create(name)
	return
}

func usrDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatalln(err)
	}
	return dir
}
