package main

import (
	"log"
	"strings"
)

type logSniffer struct {
	startSentinel map[string]bool
	startChan     chan bool
	prefix        string
	color         int
}

// Blocks until all log sentines have been encountered.
func (l *logSniffer) waitForSentinels() {
	if l.startChan == nil {
		return
	}
	<-l.startChan
}

func (l *logSniffer) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if strings.TrimRight(line, "\r\n") == "" {
			continue
		}
		log.Printf("%s\t%s", colorize(l.prefix, l.color), line)

		// Check if we need to track output messages
		if l.startChan != nil {
			l.checkSentinels(line)
		}
	}
	return len(p), nil
}

// Checks all sentinels that we expect and sends a signal
// to the start channel once all sentinels have been received.
func (l *logSniffer) checkSentinels(line string) {
	for k := range l.startSentinel {
		if strings.Index(line, k) < 0 {
			continue
		}
		delete(l.startSentinel, k)
		if len(l.startSentinel) == 0 {
			l.startChan <- true
			l.startChan = nil
		}
	}
}

func makeSniffer(prefix string, startSentinel []string) *logSniffer {
	l := &logSniffer{
		prefix: prefix,
		color:  nextColor(),
	}

	// If this twerk defines a list of output messages to wait for,
	// copy them info a set for the sniffer.
	if len(startSentinel) != 0 {
		l.startSentinel = make(map[string]bool)
		for _, line := range startSentinel {
			l.startSentinel[line] = false
		}
		l.startChan = make(chan bool)
	}
	return l
}
