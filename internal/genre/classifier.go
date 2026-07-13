package genre

import "os"

type Result struct {
	Genre      Genre
	Confidence float64
}

func Classify(path string) Result {
	m, err := LoadMetadata(path)
	if err == nil && m.Genre != Unknown {
		return Result{Genre: m.Genre, Confidence: 1.0}
	}

	entries := listDir(path)
	hasData := contains(entries, "data")
	hasInfo := contains(entries, "info")
	hasKnowl := contains(entries, "knowl")
	hasWisdom := contains(entries, "wisdom")

	if hasData && hasInfo && hasKnowl && hasWisdom {
		return Result{Genre: Manual, Confidence: 0.9}
	}

	return Result{Genre: Unknown, Confidence: 0}
}

func listDir(path string) []string {
	var names []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
