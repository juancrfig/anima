package cmd

import (
	"os"
	"fmt"
	"io"
)

func initialGreeting(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	fmt.Fprintf(w, "Hello! I am Anima.")
}
