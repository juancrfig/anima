package cli

import (
	"bytes"
    "testing"

	"github.com/stretchr/testify/assert"
    // "github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	buffer := &bytes.Buffer{}
    // Tell rootCmd to write its output to our buffer
	rootCmd.SetOut(buffer)
    rootCmd.Execute()
    output := buffer.String()

	assert.Equal(t, "Hello, Journal!\n", output)
}
