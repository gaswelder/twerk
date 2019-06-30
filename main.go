package main

import (
	"fmt"
	"os"
)

type twerkable interface {
	start(name string, t twerks) error
}

type twerks map[string]twerkable

func (tt twerks) run(name string) error {
	tt.start(name)
	select {}
}

func (tt twerks) start(name string) error {
	return tt[name].start(name, tt)
}

func main() {
	cfg, err := parseConfig("twerks.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = cfg.run("default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
