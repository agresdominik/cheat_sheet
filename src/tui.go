package main

import (
	"fmt"
	"os"
	//"sort"
	//"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type view int

type model struct {
	list          list.Model
	commands      CmdList
	currentView   view
	currentKey    string
	selectedCmd   string
	categoryIndex int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			selected, ok := m.list.SelectedItem().(item)
			if ok {
				if m.currentView == viewCategories {
					m.categoryIndex = m.list.Index()
					m.currentKey = selected.title
					m.currentView = viewCommands
					m.list.Title = m.currentKey
					m.list.SetItems(cmdItemsToList(m.commands.Get(m.currentKey)))
					m.list.Select(0)
				} else if m.currentView == viewCommands {
					m.selectedCmd = selected.title
					return m, tea.Quit
				}
			}

		case "b":

			if m.currentView == viewCommands {
				m.currentView = viewCategories
				m.list.Title = "Choose a list of commands"
				m.list.SetItems(cmdListKeysToList(m.commands))
				m.list.Select(m.categoryIndex)
			} else if m.currentView == viewCategories {
				return m, tea.Quit
			}

		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

var docStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Background(lipgloss.Color("235"))

const (
	viewCategories view = iota
	viewCommands
)

func StartTui(commands CmdList) {

	delegate := newStyledDelegate()

	items := cmdListKeysToList(commands)
	l := list.New(items, delegate, 0, 0)
	l.Title = "Choose a list of commands"
	l.Styles = newListStyles()

	m := model{
		list:        l,
		commands:    commands,
		currentView: viewCategories,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok && m.selectedCmd != "" {
		fmt.Println(m.selectedCmd)
	}
}

func cmdListKeysToList(cmds CmdList) []list.Item {
	items := []list.Item{}
	for _, group := range cmds {
		listItem := item{
			title:       group.Category,
			description: "",
		}
		items = append(items, listItem)
	}
	return items
}

func cmdItemsToList(cmds []CmdItem) []list.Item {

	items := []list.Item{}

	for _, command := range cmds {
		listItem := item{
			title:       command.CommandName,
			description: command.CommandDescription,
		}
		items = append(items, listItem)
	}
	return items
}

func newStyledDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()

	backKey := key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back/exit"),
	)
	delegate.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{backKey}
	}
	delegate.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{
			{backKey},
		}
	}

	delegate.SetSpacing(0)

	styles := list.NewDefaultItemStyles()
	styles.NormalTitle = styles.NormalTitle.
		Foreground(lipgloss.Color("110")).
		PaddingLeft(1)
	styles.NormalDesc = styles.NormalDesc.
		Foreground(lipgloss.Color("247")).
		PaddingLeft(2)
	styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("187")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("186")).
		Padding(0, 1)
	styles.SelectedDesc = styles.SelectedTitle.
		Foreground(lipgloss.Color("232")).
		Background(lipgloss.Color("230")).
		Padding(0, 1)
	styles.FilterMatch = styles.FilterMatch.
		Foreground(lipgloss.Color("226"))

	delegate.Styles = styles

	return delegate
}

func newListStyles() list.Styles {
	styles := list.DefaultStyles()

	styles.Title = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("180")).
		Background(lipgloss.Color("23")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 2)
	styles.TitleBar = styles.TitleBar.MarginBottom(1)

	styles.StatusBar = styles.StatusBar.
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252")).
		Padding(0, 1)
	styles.StatusEmpty = styles.StatusEmpty.Foreground(lipgloss.Color("60"))

	styles.PaginationStyle = styles.PaginationStyle.
		Foreground(lipgloss.Color("244")).
		PaddingLeft(4)
	styles.HelpStyle = styles.HelpStyle.
		MarginLeft(1).
		Foreground(lipgloss.Color("244"))

	styles.ActivePaginationDot = styles.ActivePaginationDot.
		Foreground(lipgloss.Color("220"))
	styles.InactivePaginationDot = styles.InactivePaginationDot.
		Foreground(lipgloss.Color("238"))

	return styles
}
