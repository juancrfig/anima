package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {

	var buffer bytes.Buffer
	Greet(&buffer)

	got := buffer.String()
	want := "Hello!\n"

	if got != want {
		t.Errorf("\nGot: %s Want: %s", got, want)
	}
}

func TestParseArgs(t *testing.T) {
	t.Run("without arguments", func(t *testing.T) {
		args := []string{"program"}

		got := ParseArgs(args)
		want := ""

		if got != want {
			t.Errorf("\nGot: %s Want: %s", got, want)
		}
	})
	t.Run("with today", func(t *testing.T) {
		args := []string{"program", "today"}

		got := ParseArgs(args)
		want := "today"

		if got != want {
			t.Errorf("\nGot: %s Want: %s", got, want)
		}
	})
}

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want string
	}{
		{"no command", "", "Hello!\n"},
		{"init command", "today", "Opening today's journal entry...\n"},
		{"unknown command", "foo", usage + "\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			RunCommand(tc.cmd, &buffer)

			got := buffer.String()
			want := tc.want

			if got != want {
				t.Errorf("\nGot: %s Want: %s", got, want)
			}
		})
	}
}
