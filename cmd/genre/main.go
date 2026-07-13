package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: genre <command> [flags]")
		fmt.Fprintln(os.Stderr, "commands: detect, validate, transcode")
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "detect":
		fs := flag.NewFlagSet("detect", flag.ExitOnError)
		if err := runDetect(fs, args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "validate":
		fs := flag.NewFlagSet("validate", flag.ExitOnError)
		if err := runValidate(fs, args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "transcode":
		fs := flag.NewFlagSet("transcode", flag.ExitOnError)
		if err := runTranscode(fs, args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
