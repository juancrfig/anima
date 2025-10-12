package main

import "testing"

func TestGreeting(t *testing.T) {
	expected := "Hello, TDD!"
	actual := Greeting()
	if actual != expected {
		t.Errorf("expected %q but got %q", expected, actual)
	}
}
