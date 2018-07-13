package parser

import (
	"log"
	"os"
)

// Mover Mover file from path to path.
type Mover interface {
	Move(from, to string) error
	Mkdir(path string) error
}

// BasicMover ...
type BasicMover struct{}

// Move ...
func (m BasicMover) Move(from, to string) error {
	return os.Rename(from, to)
}

// Mkdir ...
func (m BasicMover) Mkdir(path string) error {
	return os.Mkdir(path, 0775)
}

// FakeMover ...
type FakeMover struct{}

// Move ...
func (m FakeMover) Move(from, to string) error {
	log.Printf("fake moving\n\tfrom: %s\n\tto:   %s", from, to)
	return nil
}

// Mkdir ...
func (m FakeMover) Mkdir(path string) error {
	log.Printf("fake mkdir: %s", path)
	return nil
}
