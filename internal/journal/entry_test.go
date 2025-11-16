package journal

import (
	"testing"
	"os"
	"strings"
	"fmt"
	"path/filepath"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

func readMetadataFromFile(absPath string) (Metadata, error) {
    var meta Metadata

    data, err := os.ReadFile(absPath)
    if err != nil {
        return meta, err
    }

    parts := strings.SplitN(string(data), "---", 3)
    if len(parts) < 3 {
        return meta, fmt.Errorf("frontmatter not found")
    }

    err = yaml.Unmarshal([]byte(parts[1]), &meta)
    return meta, err
}

func TestCreateEntry(t *testing.T) {
	tmpDir := t.TempDir()
	entryPath := filepath.Join(tmpDir, "fooEntry.md")

	require.Nil(t, createEntry(entryPath))
}
