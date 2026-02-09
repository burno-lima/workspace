package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	str "workspace/pkg/strings"
	"workspace/pkg/tui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(str.AppUsage)
		fmt.Println(str.AppUsageExample)
		os.Exit(1)
	}

	basePath := os.Args[1]

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		fmt.Printf(str.ErrorPathNotExist.String()+"\n", basePath)
		os.Exit(1)
	}

	m := tui.New(basePath)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf(str.ErrorPrefix.String()+"\n", err)
		os.Exit(1)
	}
}
