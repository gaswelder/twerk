package main

import (
	"encoding/json"
	"errors"
	"log"
)

// A group of twerks that is started in a specific sequence.
type compositeTwerk struct {
	Compose [][]string `json:"compose"`
	Desc    string     `json:"desc"`
}

func parseComposite(data json.RawMessage) (*compositeTwerk, error) {
	c := &compositeTwerk{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	if len(c.Compose) == 0 {
		return nil, errors.New("empty compose field")
	}
	return c, validateJSONKeys(data, []string{"compose", "desc"})
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
