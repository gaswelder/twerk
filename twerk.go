package main

import (
	"log"
	"os/exec"
	"strings"
)

type twerk struct {
	Dir         string `json:"dir"`
	Cmd         string `json:"cmd"`
	InitMessage string `json:"initMessage"`
}

func (t *twerk) start(name string, tt twerks) error {
	cmd := exec.Command("sh", "-c", t.Cmd)
	cmd.Dir = t.Dir

	// Copy all output to stdout with prefix
	logPrefix := name
	redir := makeSniffer(logPrefix, t.InitMessage)
	cmd.Stderr = redir
	cmd.Stdout = redir

	// Start the process
	go func() {
		err := cmd.Run()
		log.Printf("%v exited: %v", name, err)
	}()

	// Wait for the sentinel on the channel
	redir.wait()
	log.Printf("-------------------- " + name + " started ------------------------")
	return nil
}

type logSniffer struct {
	startSentinel string
	startChan     chan bool
	prefix        string
}

func (l *logSniffer) wait() {
	if l.startChan == nil {
		return
	}
	<-l.startChan
}

func (l *logSniffer) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		log.Printf("%s\t%s", l.prefix, line)
		if l.startChan != nil && strings.Index(line, l.startSentinel) >= 0 {
			l.startChan <- true
			l.startChan = nil
		}
	}
	return len(p), nil
}

func makeSniffer(prefix string, startSentinel string) *logSniffer {
	l := &logSniffer{
		prefix: prefix,
	}
	if startSentinel != "" {
		l.startSentinel = startSentinel
		l.startChan = make(chan bool)
	}
	return l
}
