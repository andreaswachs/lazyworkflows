package tui

// TODO: UI is unresponsive.. check out https://github.com/charmbracelet/bubbletea/blob/79c76c680b1a6bae9cd9bc918c1d8eb336ee4ceb/examples/list-fancy/main.go
// to see what we're not doing right
import (
	"math"
	"strings"

	"github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/andreaswachs/lazyworkflows/consumer"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tabState uint8

const (
	overview tabState = iota
	workflow
)

type model struct {
	conf        appconfig.AppConfig
	selectedTab tabState
	cursorPos   map[tabState]int
	fullTable   table.Model
}

// InitialModel returns an inital model to bootstrap the UI
func InitialModel(appconfig appconfig.AppConfig) model {
	initStyles()

	cursorPos := make(map[tabState]int)
	cursorPos[overview] = 0
	cursorPos[workflow] = 0

	// Load the table data
	// TODO: the table needs to be updated between frames
	// as to make the width of the last column dynamic
	columns := []table.Column{
		{Title: "Owner", Width: 10},
		{Title: "Repo", Width: 20},
		{Title: "Name", Width: 80}, // This width needs to be dynamic
	}

	rows := []table.Row{}
	api := consumer.New()

	// This should not be the responsibility of the UI to handle workflows
	for _, repo := range appconfig.Repos {
		workflows, err := api.List(repo)
		if err != nil {
			panic(err)
		}

		for _, workflow := range workflows {
			rows = append(rows, table.Row{repo.Owner, repo.Repo, workflow.Name})
		}
	}

	fullTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	fullTable.SetStyles(tableStyle)

	return model{conf: appconfig, selectedTab: workflow, cursorPos: cursorPos, fullTable: fullTable}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		width = msg.Width
		m.fullTable.SetWidth(msg.Width - 2)
		return m, nil
	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			m.fullTable.MoveDown(1)
			return m, nil
		case "k", "up":
			m.fullTable.MoveUp(1)
			return m, nil
		case "h", "left":
			m.selectedTab = previousTab(m.selectedTab)
			return m, nil
		case "l", "right":
			m.selectedTab = nextTab(m.selectedTab)
			return m, nil
		}

		// Return the updated model to the Bubble Tea runtime for processing.
		// Note that we're not returning a command.
	}
	return m, nil
}

func (m model) View() string {
	builder := strings.Builder{}

	renderTabs(&builder, &m)
	builder.WriteString("\n")
	renderBody(&builder, &m)

	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("Press q or ctrl+c to quit")

	return builder.String()
}

func renderTabs(builder *strings.Builder, m *model) {
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderSingleTab(overview, m.selectedTab),
		renderSingleTab(workflow, m.selectedTab),
	)
	gap := tabGap.Render(strings.Repeat(" ", int(math.Abs(float64(width-len(row)-2)))))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	builder.WriteString(row)
	builder.WriteString("\n")
}

func renderSingleTab(currentTab tabState, selectedTab tabState) string {
	if currentTab == selectedTab {
		return activeTab.Render(tabStateToTab(currentTab))
	}
	return tab.Render(tabStateToTab(currentTab))
}

func renderBody(builder *strings.Builder, m *model) {
	switch m.selectedTab {
	case overview:
		renderOverview(builder, m)
	case workflow:
		renderWorkflow(builder, m)
	}
}

func renderOverview(builder *strings.Builder, m *model) {
	builder.WriteString(baseStyle.Render(m.fullTable.View()))
}

func renderWorkflow(builder *strings.Builder, m *model) {
	builder.WriteString("TODO\n")
}

func tabStateToTab(selectedTab tabState) string {
	switch selectedTab {
	case overview:
		return "Overview"
	case workflow:
		return "Workflow"
	}
	return ""
}

func nextTab(selectedTab tabState) tabState {
	switch selectedTab {
	case overview:
		return workflow
	case workflow:
		return overview
	}
	return overview
}

func previousTab(selectedTab tabState) tabState {
	switch selectedTab {
	case overview:
		return workflow
	case workflow:
		return overview
	}
	return overview
}

func moveCursorDown(m *model) {
	if m.selectedTab == overview {
		if m.cursorPos[m.selectedTab] < len(m.conf.Repos)-1 {
			m.cursorPos[m.selectedTab]++
		} else {
			m.cursorPos[m.selectedTab] = 0
		}
	} else {
		// TODO: Implement cursor movement for other tabs
		m.cursorPos[m.selectedTab]++
	}
}

func moveCursorUp(m *model) {
	if m.selectedTab == overview {
		if m.cursorPos[m.selectedTab] > 0 {
			m.cursorPos[m.selectedTab]--
		} else {
			m.cursorPos[m.selectedTab] = len(m.conf.Repos) - 1
		}
	} else {
		// TODO: Implement bounds check for other tabs
		m.cursorPos[m.selectedTab]--
	}
}
