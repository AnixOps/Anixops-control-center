package main

import (
	"testing"
)

func TestMainImport(t *testing.T) {
	// This test verifies the package compiles correctly
	// The main function starts an interactive TUI, so we can't test it directly
	// without mocking the terminal
}

func TestPackageStructure(t *testing.T) {
	// Verify the package is valid
	// This is a minimal test for the TUI entry point
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
}