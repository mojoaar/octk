package main

import (
	"fmt"
	"os"

	"github.com/mojoaar/octk/cmd"
	"github.com/mojoaar/octk/internal"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "help", "--help", "-h":
		cmd.PrintHelp()

	case "version", "--version", "-v":
		fmt.Printf("octk v%s\n", internal.ToolVersion)

	case "list":
		if err := cmd.List(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "backup":
		if err := cmd.Backup(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "restore":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: octk restore <file.tar.gz>")
			os.Exit(1)
		}
		if err := cmd.Restore(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		fmt.Fprintln(os.Stderr, "run 'octk help' for usage")
		os.Exit(1)
	}
}
