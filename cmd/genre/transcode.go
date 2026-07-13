package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/genre"
	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/genre/transcoder"
)

func runTranscode(fs *flag.FlagSet, args []string) error {
	input := fs.String("input", "", "path to source document directory")
	output := fs.String("output", "", "output directory")
	fs.Parse(args)

	if *input == "" {
		return fmt.Errorf("-input is required")
	}
	if *output == "" {
		return fmt.Errorf("-output is required")
	}

	m, err := genre.LoadMetadata(*input)
	if err != nil {
		return fmt.Errorf("loading genre metadata: %w", err)
	}

	t := transcoder.New()
	archivePath, err := t.ManualToArchive(*input, *output, m)
	if err != nil {
		return fmt.Errorf("transcode error: %w", err)
	}

	fmt.Fprintf(os.Stderr, "transcoded: %s → %s\n", *input, archivePath)
	return nil
}
