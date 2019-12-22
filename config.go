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

	// Parse into an opaque string->?? map first.
	kv := make(map[string]json.RawMessage)
	err = json.Unmarshal(data, &kv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", path, err)
	}

	cfg := make(twerks)
	for name, data := range kv {
		t, err := parseConfigNode(data, &cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse twerks[%v]: %v", name, err)
		}
		cfg[name] = t
	}

	return cfg, nil
}

func parseConfigNode(data json.RawMessage, tt *twerks) (twerkable, error) {
	// Is it a regular twerk?
	t, err1 := parseTwerk(data)
	if t != nil {
		return t, err1
	}

	// Is it a composite?
	c, err2 := parseComposite(data, tt)
	if c != nil {
		return c, err2
	}

	return nil, fmt.Errorf("unrecognized node format (twerk: %v, composite: %v)", err1, err2)
}

func validateJSONKeys(data json.RawMessage, keys []string) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	for key := range m {
		if !contains(keys, key) {
			return fmt.Errorf("unknown field: %s", key)
		}
	}
	return nil
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
