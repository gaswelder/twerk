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
	err := tt.start(name)
	if err != nil {
		return err
	}
	select {}
}

func (tt twerks) start(name string) error {
	if tt[name] == nil {
		return fmt.Errorf("twerk %s is not defined", name)
	}
	return tt[name].start(name, tt)
}

func main() {
	cfg, err := parseConfig("twerks.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: twerkgraf [taskname]\n")
		os.Exit(1)
	}

	taskName := "default"
	if len(os.Args) == 2 {
		taskName = os.Args[1]
	}

	err = cfg.run(taskName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
