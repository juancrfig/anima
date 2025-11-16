package cmd

import (
	"testing"
	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialGreeting(t *testing.T) {
	var buf bytes.Buffer

	InitialGreeting(&buf)

	require.NotNil(t, rootCmd)
	assert.Equal(t, "Hello! I'm Anima.", buf.String())
}

