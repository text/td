package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var logger = log.New(os.Stdout, "today: ", log.Lshortfile)

func main() {
	p := &Program{
		create:  create,
		envDir:  os.Getenv("TODAYPATH"),
		open:    open,
		pre:     "\033[1m",
		started: time.Now(),
		suf:     "\033[0m",
		usrDir:  usrDir,
	}
	flag.BoolVar(&p.printDuration, "printDuration", true, "")
	flag.BoolVar(&p.printRange, "printRange", false, "")
	flag.DurationVar(&p.roundDur, "roundDuration", time.Duration(time.Minute), "")
	flag.DurationVar(&p.truncateDur, "truncateDuration", time.Duration(0), "")
	flag.StringVar(&p.dir, "homeDir", "", "home directory, if not set $TODAYPATH or $HOME/.today is used")
	flag.StringVar(&p.prefix, "prefix", "", "")
	flag.Parse()

	if err := p.Load(); err != nil {
		logger.Fatal(err)
	}
	text := strings.Join(flag.Args(), " ")
	t := time.Now()
	if flag.NArg() == 0 {
		p.Print(t)
		return
	}
	p.Add(t, text)
	if err := p.Save(); err != nil {
		logger.Fatal(err)
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
