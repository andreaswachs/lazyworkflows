package tui

import (
	"github.com/andreaswachs/lazyworkflows/appconfig"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	conf appconfig.AppConfig
}

func InitialModel(appconfig appconfig.AppConfig) model {
	return model{conf: appconfig}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// Return the updated model to the Bubble Tea runtime for processing.
		// Note that we're not returning a command.
	}
	return m, nil
}

func (m model) View() string {
	return "Press q or ctrl+c to quit."
}
