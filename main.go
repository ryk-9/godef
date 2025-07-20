// File: main.go
// Description: The main entry point for the GoDef application.

package main

import (
	"godef/cmd"
)

// main is the entry point of the application.
// It simply executes the root command of the CLI application.
func main() {
	cmd.Execute()
}
