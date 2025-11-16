package journal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	testLocationURL := "http://ip-api.com/json/24.48.0.1?fields=status,message,country,city"

	r, err := getLocation(testLocationURL)

	assert.Nil(t, err)
	assert.Equal(t, "Montreal", r[0])
	assert.Equal(t, "Canada", r[1])
}
