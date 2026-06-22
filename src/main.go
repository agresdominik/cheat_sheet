package main

import (
	"fmt"
	"os"
	"flag"
	"encoding/json"
	"log"
	"sort"
	"path/filepath"
)

type CmdGroup struct {
    Category string
    Commands []CmdItem
}

type CmdList []CmdGroup

type CmdItem struct {
	CommandName        string `json:"command"`
	CommandDescription string `json:"desc"`
}

func (list CmdList) Get(category string) []CmdItem {
	for _, group := range list {
		if group.Category == category {
			return group.Commands
		}
	}
	return nil
}


func main() {

	configFlag := flag.String("config", "", "Specify a config file")
	helpFlag := flag.Bool("help", false, "Show help")
	newFlag := flag.Bool("new", false, "Add new command to config")

	flag.Parse()

	switch {
	case *helpFlag:
		printHelp()
		return
	case *newFlag:
		HandleInput()
		return
	}

	configFile := *configFlag
	if configFile == "" {
		configFile = defaultConfigPath()
	}

	commands, err := loadCommands(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("No config found at %s\nPass a config explicitly with --config <file>.", configFile)
		}
		log.Fatalf("Cannot load commands file: %v", err)
	}

	StartTui(commands)
}

// defaultConfigPath returns the bundled command list installed next to the
// binary (<prefix>/share/cheatsh/commands.json, replaced on upgrade). When
// running from a source checkout that file doesn't exist, so it falls back to
// data/commands.json relative to the working directory.
func defaultConfigPath() string {
	if exe, err := os.Executable(); err == nil {
		if real, err := filepath.EvalSymlinks(exe); err == nil {
			exe = real
		}
		p := filepath.Join(filepath.Dir(exe), "..", "share", "cheatsh", "commands.json")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return filepath.Join("data", "commands.json")
}

func printHelp() {
	fmt.Println(`Usage: cheatsh [options]
	Select a command to copy it to the system clipboard.
	Without --config the bundled command list is used; pass --config to use your own.
	Options:
	  --config <file>  Use a custom config file
	  --help           Show this help message
	  --new            Add a new command to the config file`)
}

func loadCommands(path string) (CmdList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string][]CmdItem
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	cmdList := make(CmdList, 0, len(raw))
	for category, commands := range raw {
		cmdList = append(cmdList, CmdGroup{
			Category: category,
			Commands: commands,
		})
	}

	sort.Slice(cmdList, func(i, j int) bool {
		return cmdList[i].Category < cmdList[j].Category
	})

	return cmdList, nil
}
