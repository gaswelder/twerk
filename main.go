package main

import (
	"fmt"
	"os"
)

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

	// If no arguments given, show the list of available commands.
	if len(os.Args) == 1 {
		help(cfg)
		os.Exit(0)
	}

	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: twerkgraf [taskname]\n")
		os.Exit(1)
	}

	taskName := os.Args[1]
	err = cfg.run(taskName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func help(tt twerks) {
	for name, t := range tt {
		d := t.desc()
		if d == "" {
			continue
		}
		fmt.Printf("	%s	%s\n", name, d)
	}
}
