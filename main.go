package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/eneslevent/spomac/internal/ui"
)

var version = "0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "Show version information")
	flag.BoolVar(showVersion, "v", false, "Show version information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "spomac - Lightweight Spotify controller for macOS terminal\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  spomac [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -v, --version    Show version information\n")
		fmt.Fprintf(os.Stderr, "  -h, --help       Show this help message\n\n")
		fmt.Fprintf(os.Stderr, "Key bindings:\n")
		fmt.Fprintf(os.Stderr, "  space, 2         Play/Pause\n")
		fmt.Fprintf(os.Stderr, "  left, 1          Previous track\n")
		fmt.Fprintf(os.Stderr, "  right, 3         Next track\n")
		fmt.Fprintf(os.Stderr, "  up               Volume up (+5)\n")
		fmt.Fprintf(os.Stderr, "  down             Volume down (-5)\n")
		fmt.Fprintf(os.Stderr, "  q, ctrl+c        Quit\n")
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("spomac v%s\n", version)
		return
	}

	p := tea.NewProgram(ui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
