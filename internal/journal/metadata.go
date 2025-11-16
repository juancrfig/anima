
package journal

import (
	"os"
	"fmt"
	"strings"
	"net/http"
	"io"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

const locationURL string = "http://ip-api.com/json/?fields=status,message,country,city"

type Location struct {
	Country string
	City    string
}

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

func getLocation(url string) ([2]string, error) {
	if url == "" {
		url = locationURL
	}
	resp, err := http.Get(url)
	if err != nil {
		return [2]string{}, err
	}
	
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return [2]string{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return [2]string{}, err
	}

	var locData Location
	err = json.Unmarshal(body, &locData)
	if err != nil {
		return [2]string{}, err
	}

	s := [2]string{locData.City, locData.Country}

	return s, nil
}
