package cmd

import (
	"testing"
	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialGreeting(t *testing.T) {
	var buf bytes.Buffer

	initialGreeting(&buf)

	require.NotNil(t, rootCmd)
	assert.Equal(t, "Hello! I am Anima.", buf.String())
}

