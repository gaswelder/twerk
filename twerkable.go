package main

// Twerkable is a task which can be run.
type twerkable interface {
	// desc returns the task's description.
	desc() string

	// start starts the task.
	start(name string, t twerks) error

	// wait blocks until the task finishes.
	wait() error
}
