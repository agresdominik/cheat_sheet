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
	title 			string
	description 	string
}
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type view int

type model struct {
	list 		list.Model
	commands 	CmdList
	currentView view
	currentKey	string
	selectedCmd string
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
							m.currentKey = selected.title
							m.currentView = viewCommands
							m.list.Title = m.currentKey
							m.list.SetItems(cmdItemsToList(m.commands.Get(m.currentKey)))
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

var docStyle = lipgloss.NewStyle().Margin(1, 2)

const (
	viewCategories view = iota
	viewCommands
)

func StartTui(commands CmdList) {

	delegate := list.NewDefaultDelegate()
	backKey := key.NewBinding(
	    key.WithKeys("b"),
	    key.WithHelp("b", "back"),
	)
	delegate.ShortHelpFunc = func() []key.Binding{
		return []key.Binding{
			backKey,
		}
	}
	delegate.FullHelpFunc = func() [][]key.Binding{
		return [][]key.Binding{
			{backKey},
		}
	}

	items := cmdListKeysToList(commands)
	l := list.New(items, delegate, 0, 0)
	l.Title = "Choose a list of commands"

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
			title: command.CommandName,
			description: command.CommandDescription,
		}
		items = append(items, listItem)
	}
	return items
}
