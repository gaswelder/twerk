package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func parseConfig(path string) (twerks, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", path, err)
	}

	// Parse into an opaque string->?? map first.
	kv := make(map[string]json.RawMessage)
	err = json.Unmarshal(data, &kv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", path, err)
	}

	cfg := make(twerks)
	for name, data := range kv {
		t, err := parseConfigNode(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse twerks[%v]: %v", name, err)
		}
		cfg[name] = t
	}

	return cfg, nil
}

func parseConfigNode(data json.RawMessage) (twerkable, error) {
	// Is it a composite?
	c := make(compositeTwerk, 0)
	err := json.Unmarshal(data, &c)
	if err == nil {
		return c, nil
	}

	// Is it a regular twerk?
	t := new(twerk)
	err = json.Unmarshal(data, &t)
	if err == nil {
		return t, nil
	}
	return nil, errors.New("unrecognized node format")
}
