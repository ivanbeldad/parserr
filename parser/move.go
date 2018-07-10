package parser

import (
	"fmt"
	"io"
	"os"
)

// Move Move file from path to path.
type Move interface {
	Move(from, to string) error
}

// BasicMove ...
type BasicMove struct{}

// Move ...
func (m BasicMove) Move(from, to string) error {
	return os.Rename(from, to)
}

// DiskMove Move file between disks
type DiskMove struct{}

// Move ...
func (m DiskMove) Move(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

// FakeMove ...
type FakeMove struct{}

// Move ...
func (m FakeMove) Move(from, to string) error {
	fmt.Printf("moveing from %s to %s", from, to)
	return nil
}
