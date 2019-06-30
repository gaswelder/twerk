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
	// Is it a regular twerk?
	t := new(twerk)
	err := json.Unmarshal(data, &t)
	if err == nil && t.Cmd != "" {
		err = validateJSONKeys(data, []string{"cmd", "desc", "dir", "logPrefix", "initMessages", "env"})
		return t, err
	}

	// Is it a composite?
	c := newComposite()
	err = json.Unmarshal(data, &c)
	if err == nil && len(c.Compose) != 0 {
		err = validateJSONKeys(data, []string{"compose", "desc"})
		return c, err
	}

	return nil, errors.New("unrecognized node format")
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
