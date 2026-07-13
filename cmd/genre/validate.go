package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/genre"
	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/genre/validator"
)

func runValidate(fs *flag.FlagSet, args []string) error {
	input := fs.String("input", "", "path to document directory")
	fs.Parse(args)

	if *input == "" {
		return fmt.Errorf("-input is required")
	}

	m, err := genre.LoadMetadata(*input)
	if err != nil {
		return fmt.Errorf("loading genre metadata: %w", err)
	}

	rules := validator.ForGenre(m.Genre)
	report := rules.Validate(*input)
	report.Genre = m.Genre

	enc := NewReportEncoder(os.Stdout)
	return enc.Encode(report)
}

type ReportEncoder struct {
	w   interface{ Write([]byte) (int, error) }
	fmt bool
}

func NewReportEncoder(w interface{ Write([]byte) (int, error) }) *ReportEncoder {
	return &ReportEncoder{w: w}
}

func (e *ReportEncoder) Encode(r *validator.Report) error {
	_, err := fmt.Fprintf(e.w, "Genre: %s\n", r.Genre)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(e.w, "Passed: %d, Warnings: %d, Errors: %d\n", r.Passed, r.Warnings, r.Errors)
	if err != nil {
		return err
	}
	for _, item := range r.Items {
		_, err = fmt.Fprintf(e.w, "  [%s] %s\n", item.Level, item.Message)
		if err != nil {
			return err
		}
	}
	return nil
}
