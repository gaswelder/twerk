package main

import "log"

// A group of twerks that is started in a specific sequence.
type compositeTwerk [][]string

func (t compositeTwerk) start(name string, tt twerks) error {
	for _, group := range t {
		for _, name := range group {
			err := tt.start(name)
			if err != nil {
				log.Fatalf("failed to start %s: %v", name, err)
			}
		}
	}

	return nil
}
