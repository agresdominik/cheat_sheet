package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Individual Command Items Later Extracted
type CmdItem struct {
	CommandName        string `json:"command"`
	CommandDescription string `json:"desc"`
}

// All commands
type CmdMap map[string][]CmdItem

// Read Json-File and return it
func loadCommands(path string) (CmdMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m CmdMap
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Function Called when a subset of functions has been selected
func showCommandList(app *tview.Application, lang string, cmds []CmdItem, langList *tview.List) {
	cmdList := tview.NewList()
	cmdList.SetBorder(true).SetTitle(fmt.Sprintf("%s commands", lang))

	for _, c := range cmds {
		c := c
		cmdList.AddItem(c.CommandName, c.CommandDescription, 0, func() {
			//fmt.Println(c.CommandName)
			app.Suspend(func() {
				fmt.Println(c.CommandName) // this now prints visibly
			})
			app.Stop()
			os.Exit(0)
		})
	}

	cmdList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'b' || event.Key() == tcell.KeyEscape {
			// Restore the language list
			app.SetRoot(langList, true)
			return nil
		}
		return event
	})

	app.SetRoot(cmdList, true)
}

// Main Function
// Initialises the first page and starts app
func main() {
	commands, err := loadCommands("data/commands.json")
	if err != nil {
		log.Fatalf("Cannot load commands file: %v", err)
	}

	app := tview.NewApplication()

	// Sort alphabetically
	langs := make([]string, 0, len(commands))
	for lang := range commands {
		langs = append(langs, lang)
	}
	sort.Strings(langs)

	langList := tview.NewList()
	langList.SetBorder(true).SetTitle("Select a topic")

	for _, lang := range langs {
		lang := lang // capture loop variable
		langList.AddItem(lang, "", 0, func() {
			showCommandList(app, lang, commands[lang], langList)
		})
	}

	if err := app.SetRoot(langList, true).Run(); err != nil {
		log.Fatalf("Error initialising TUI: %v", err)
	}
}
