package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
)

type twerk struct {
	Dir          string            `json:"dir"`
	Cmd          string            `json:"cmd"`
	InitMessages []string          `json:"initMessages"`
	Env          map[string]string `json:"env"`
	Desc         string            `json:"desc"`

	end chan error
}

func parseTwerk(data json.RawMessage) (*twerk, error) {
	t := new(twerk)
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	if t.Cmd == "" {
		return nil, errors.New("not a twerk: cmd field is missing")
	}
	return t, validateJSONKeys(data, []string{"cmd", "desc", "dir", "logPrefix", "initMessages", "env"})
}

// Converts an env var map to a list of "name=value" pairs.
func envList(m map[string]string) []string {
	list := make([]string, len(m))
	for name, val := range m {
		list = append(list, name+"="+val)
	}
	return list
}

func (t *twerk) start(name string, tt twerks) error {
	log.Printf("-------------------- starting " + name + " ------------------------")
	if t.end != nil {
		return errors.New("end channel already exists")
	}
	t.end = make(chan error, 0)

	cmd := exec.Command("sh", "-c", t.Cmd)
	cmd.Env = append(os.Environ(), envList(t.Env)...)
	cmd.Dir = t.Dir

	// Copy all output to stdout with prefix
	logPrefix := name
	logger := makeSniffer(logPrefix, t.InitMessages)
	cmd.Stderr = logger
	cmd.Stdout = logger

	// Start the process
	go func() {
		err := cmd.Run()
		log.Printf("%v exited: %v", name, err)
		t.end <- err
		t.end = nil
	}()

	// Wait for the sentinel on the channel
	logger.waitForSentinels()
	log.Printf("-------------------- " + name + " started ------------------------")
	return nil
}

func (t *twerk) desc() string {
	return t.Desc
}

func (t *twerk) wait() error {
	return <-t.end
}
