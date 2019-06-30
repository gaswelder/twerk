package main

import "log"

// A group of twerks that is started in a specific sequence.
type compositeTwerk struct {
	Compose [][]string `json:"compose"`
	Desc    string     `json:"desc"`
}

func newComposite() *compositeTwerk {
	return &compositeTwerk{}
}

func (t compositeTwerk) start(name string, tt twerks) error {
	for _, group := range t.Compose {
		for _, name := range group {
			err := tt.start(name)
			if err != nil {
				log.Fatalf("failed to start %s: %v", name, err)
			}
		}
	}

	return nil
}

func (t compositeTwerk) desc() string {
	return t.Desc
}
