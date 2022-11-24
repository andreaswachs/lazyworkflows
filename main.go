package main

import (
	"fmt"
	"os"

	appConfig "github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/andreaswachs/lazyworkflows/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := appConfig.New()

	err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config file. See error msg.")
		os.Exit(0)
	}

	p := tea.NewProgram(tui.InitialModel(*config), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not start program. See error msg.")
		os.Exit(0)
	}

}
