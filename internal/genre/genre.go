package genre

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Genre string

const (
	Unknown Genre = "unknown"
	Manual  Genre = "Manual"
	Archive Genre = "Archive"
	Journal Genre = "Journal"
	Context Genre = "Context"
)

type Contract struct {
	Form             string   `json:"form"`
	Granularity      string   `json:"granularity"`
	LifecycleStages  []string `json:"lifecycle_stages"`
}

type Metadata struct {
	Genre    Genre    `json:"genre"`
	Contract Contract `json:"contract"`
}

func LoadMetadata(path string) (*Metadata, error) {
	raw, err := os.ReadFile(filepath.Join(path, "genre.json"))
	if err != nil {
		return nil, err
	}
	var m Metadata
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
