package tui

// TODO: UI is unresponsive.. check out https://github.com/charmbracelet/bubbletea/blob/79c76c680b1a6bae9cd9bc918c1d8eb336ee4ceb/examples/list-fancy/main.go
// to see what we're not doing right
import (
	"fmt"
	"math"
	"strings"

	"github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/andreaswachs/lazyworkflows/consumer"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tabState uint8

const (
	overview tabState = iota
	workflow
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

type model struct {
	conf         appconfig.AppConfig
	selectedTab  tabState
	cursorPos    map[tabState]int
	keys         *listKeyMap
	workflowList list.Model
	delegateKeys *delegateKeyMap
}

// InitialModel returns an inital model to bootstrap the UI
func InitialModel(appconfig appconfig.AppConfig) model {
	cursorPos := make(map[tabState]int)
	cursorPos[overview] = 0
	cursorPos[workflow] = 0

	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	delegate := newItemDelegate(delegateKeys)

	// TODO: this needs to be moved somewhere else!!!
	apiConsumer := consumer.New()
	var itemList []list.Item

	for _, repo := range appconfig.Repos {
		workflows, err := apiConsumer.List(repo)
		if err != nil {
			panic(err)
		}

		for _, workflow := range workflows {
			itemList = append(itemList, item{
				title:       workflow.Name,
				description: workflow.Path,
			})
		}
	}
	workflowList := list.New(itemList, delegate, 0, 0)
	workflowList.Title = "Workflows"
	workflowList.Styles.Title = titleStyle
	workflowList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
			listKeys.insertItem,
		}
	}
	return model{workflowList: workflowList, keys: listKeys, delegateKeys: delegateKeys, conf: appconfig, selectedTab: workflow, cursorPos: cursorPos}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		width = msg.Width
		h, v := appStyle.GetFrameSize()
		m.workflowList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	// Is it a key press?
	case tea.KeyMsg:
		if m.workflowList.FilterState() == list.Filtering {
			break
		}
		// Cool, what was the actual key pressed?

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.workflowList.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.workflowList.ShowTitle()
			m.workflowList.SetShowTitle(v)
			m.workflowList.SetShowFilter(v)
			m.workflowList.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.workflowList.SetShowStatusBar(!m.workflowList.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.workflowList.SetShowHelp(!m.workflowList.ShowHelp())
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
	var listContent []string
	for i, repo := range m.conf.Repos {
		text := fmt.Sprintf("%s/%s", repo.Owner, repo.Repo)
		var label string
		if i == m.cursorPos[overview] {
			label = listSelected(text)
		} else {
			label = listItem(text)
		}
		listContent = append(listContent, label)
	}

	list := lipgloss.JoinVertical(lipgloss.Top, listContent...)
	builder.WriteString(list)
}

func renderWorkflow(builder *strings.Builder, m *model) {
	builder.WriteString(m.workflowList.View())
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

// Credits: https://github.com/charmbracelet/bubbletea/blob/master/examples/list-fancy/main.go
func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}
