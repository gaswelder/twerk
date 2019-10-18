package main

import "fmt"

type twerks map[string]twerkable

func (tt twerks) run(name string) error {
	err := tt.start(name)
	if err != nil {
		return err
	}
	return tt.wait(name)
}

func (tt twerks) start(name string) error {
	if tt[name] == nil {
		return fmt.Errorf("twerk %s is not defined", name)
	}
	return tt[name].start(name, tt)
}

func (tt twerks) wait(name string) error {
	return tt[name].wait()
}
