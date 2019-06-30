package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func parseConfig(path string) (twerks, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", path, err)
	}

	kv := make(map[string]json.RawMessage)
	err = json.Unmarshal(data, &kv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", path, err)
	}

	cfg := make(twerks)
	for name, data := range kv {
		c := make(compositeTwerk, 0)
		t := new(twerk)
		err := json.Unmarshal(data, &c)
		if err == nil {
			cfg[name] = c
			continue
		}
		err = json.Unmarshal(data, &t)
		if err == nil {
			cfg[name] = t
			continue
		}
		return nil, fmt.Errorf("failed to parse twerks[%v]: %v", name, err)
	}

	return cfg, nil
}
