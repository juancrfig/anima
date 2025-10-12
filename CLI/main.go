package main

import (
	"os"
)

func main() {
	cmd := ParseArgs(os.Args)
	RunCommand(cmd, os.Stdout)
}
