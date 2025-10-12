package main

import (
	"fmt"
	"io"
)

var usage string = `
Welcome to Anima CLI
Available arguments:
- today: Open today's journal entry`

func Greet(w io.Writer) {
	fmt.Fprintln(w, "Hello!")
}

func ParseArgs(a []string) string {
	if len(a) > 1 {
		return a[1]
	}
	return ""
}

func RunCommand(cmd string, w io.Writer) {
	switch cmd {
	case "":
		Greet(w)
	case "today":
		fmt.Fprintln(w, "Opening today's journal entry...")
	default:
		fmt.Fprintln(w, usage)
	}

}
