package journal

import (
	"testing"
	"io"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestDetectFrontmatter(t *testing.T) {
	cases := []struct{
		name   string
		input  io.Reader
		want   bool
	}{
		{
			name: "Empty file",
			input: strings.NewReader(""),
			want:  false,
		},
		{
			name: "Simple Text",
			input: strings.NewReader("Not what you look for"),
			want: false,
		},
		{
			name: "Frontmatter",
			input: strings.NewReader("---"),
			want: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := DetectFrontmatter(c.input)

			assert.Nil(t, err)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestGetLocation(t *testing.T) {
	testLocationURL := "http://ip-api.com/json/24.48.0.1?fields=status,message,country,city"

	r, err := getLocation(testLocationURL)

	assert.Nil(t, err)
	assert.Equal(t, "Montreal", r[0])
	assert.Equal(t, "Canada", r[1])
}
