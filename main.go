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

	case "version", "--version":
		fmt.Printf("octk v%s\n", internal.ToolVersion)

	case "list":
		jsonOut := hasFlag("--json")
		verbose := hasFlag("-v")
		if err := cmd.List(jsonOut, verbose); err != nil {
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
		if os.Args[1] == "-v" {
			fmt.Printf("octk v%s\n", internal.ToolVersion)
			return
		}
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		fmt.Fprintln(os.Stderr, "run 'octk help' for usage")
		os.Exit(1)
	}
}

func hasFlag(flag string) bool {
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == flag {
			return true
		}
	}
	return false
}
