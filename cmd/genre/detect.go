package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/genre"
)

func runDetect(fs *flag.FlagSet, args []string) error {
	input := fs.String("input", "", "path to document directory")
	fs.Parse(args)

	if *input == "" {
		return fmt.Errorf("-input is required")
	}

	result := genre.Classify(*input)
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
	return nil
}
