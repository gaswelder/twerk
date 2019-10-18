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

	twerksTable *twerks
	end         chan error
}

func parseComposite(data json.RawMessage, tt *twerks) (*compositeTwerk, error) {
	c := &compositeTwerk{
		twerksTable: tt,
	}
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	if len(c.Compose) == 0 {
		return nil, errors.New("empty compose field")
	}
	return c, validateJSONKeys(data, []string{"compose", "desc"})
}

func (t compositeTwerk) start(name string) error {
	if t.end != nil {
		return errors.New("non-nil end channel")
	}

	subtwerks := make([]string, 0)
	for _, group := range t.Compose {
		for _, name := range group {
			err := t.twerksTable.start(name)
			if err != nil {
				log.Fatalf("failed to start %s: %v", name, err)
			}
			subtwerks = append(subtwerks, name)
		}
	}

	t.end = make(chan error)
	go func() {
		var err error
		for _, name := range subtwerks {
			err = t.twerksTable.wait(name)
			if err != nil {
				break
			}
		}
		t.end <- err
		t.end = nil
	}()

	return nil
}

func (t compositeTwerk) desc() string {
	return t.Desc
}

func (t compositeTwerk) wait() error {
	return <-t.end
}
