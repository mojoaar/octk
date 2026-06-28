package cmd

import (
	"fmt"
	"strings"
)

const innerWidth = 41

var logoLines = []string{
	"   ██████╗  ██████╗████████╗██╗  ██╗",
	"  ██╔═══██╗██╔════╝╚══██╔══╝██║ ██╔╝",
	"  ██║   ██║██║        ██║   █████╔╝",
	"  ██║   ██║██║        ██║   ██╔═██╗",
	"  ╚██████╔╝╚██████╗   ██║   ██║  ██╗",
	"   ╚═════╝  ╚═════╝   ╚═╝   ╚═╝  ╚═╝",
}

func PrintLogo() {
	horiz := strings.Repeat("─", innerWidth)
	fmt.Println("╭" + horiz + "╮")
	fmt.Println(frameline(""))

	for _, line := range logoLines {
		fmt.Println(frameline(line))
	}

	fmt.Println(frameline(""))
	fmt.Println(frameline(center("OpenCode ToolKit", innerWidth)))
	fmt.Println(frameline(""))
	fmt.Println("╰" + horiz + "╯")
}

func PrintHelp() {
	PrintLogo()
	fmt.Println()
	fmt.Println("Backup and restore your OpenCode skills, rules, agents, and config.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  octk <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help              Show this help text")
	fmt.Println("  version           Show version information")
	fmt.Println("  list              List all discoverable skills [--json] [-v]")
	fmt.Println("  backup            Backup all skills to a tar.gz archive")
	fmt.Println("  restore <file>    Restore skills from a tar.gz archive")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  octk list --json")
	fmt.Println("  octk list -v")
	fmt.Println("  octk backup")
	fmt.Println("  octk restore octk-skills-2026-06-28.tar.gz")
}

func frameline(content string) string {
	dw := displayWidth(content)
	pad := innerWidth - dw
	if pad < 0 {
		pad = 0
	}
	return "│" + content + strings.Repeat(" ", pad) + "│"
}

func center(text string, width int) string {
	dw := displayWidth(text)
	if dw >= width {
		return text
	}
	pad := (width - dw) / 2
	return strings.Repeat(" ", pad) + text
}
