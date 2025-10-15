package main

import (
	"fmt"
	"sort"
	"log"
	"os"
	"encoding/json"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var commandsFile []byte

type CmdItem struct {
	CommandName        string `json:"command"`
	CommandDescription string `json:"desc"`
}
type CmdMap map[string][]CmdItem


func StartTui(configFilePath ...string) {

	var configFilePathString string
    if len(configFilePath) > 0 {
        configFilePathString = configFilePath[0]
    } else {
        configFilePathString = "/etc/cheatsh/commands.json"
    }
	commands, err := loadCommands(configFilePathString)
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
		lang := lang
		langList.AddItem(lang, "", 0, func() {
			showCommandList(app, lang, commands[lang], langList)
		})
	}

	langList.SetWrapAround(false).SetHighlightFullLine(true)
	langList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'b' || event.Key() == tcell.KeyEscape {
			app.Stop()
			os.Exit(0)
		}
		return event
	})

	if err := app.SetRoot(langList, true).Run(); err != nil {
		log.Fatalf("Error initialising TUI: %v", err)
	}
}

func showCommandList(app *tview.Application, lang string, cmds []CmdItem, langList *tview.List) {
	cmdList := tview.NewList()
	cmdList.SetBorder(true).SetTitle(fmt.Sprintf("%s commands", lang))

	for _, c := range cmds {
		c := c
		cmdList.AddItem(c.CommandName, c.CommandDescription, 0, func() {
			app.Suspend(func() {
				fmt.Println(c.CommandName)
			})
			app.Stop()
			os.Exit(0)
		})
	}

	cmdList.SetWrapAround(false).SetHighlightFullLine(true)
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

func loadCommands(path string) (CmdMap, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m CmdMap
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
