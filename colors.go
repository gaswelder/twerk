package main

import (
	"fmt"
	"sync"
)

// ANSI terminal color numbers.
// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
const (
	black = 30 + iota
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

func colorize(s string, color int) string {
	return fmt.Sprintf("\033[%dm%s\033[%dm", color, s, white)
}

var m sync.Mutex
var colorIndex int

var colors = []int{
	red,
	green,
	yellow,
	blue,
	magenta,
	cyan,
}

func nextColor() int {
	m.Lock()
	defer m.Unlock()

	c := colors[colorIndex]
	colorIndex++
	if colorIndex >= len(colors) {
		colorIndex = 0
	}
	return c
}
