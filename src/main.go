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

	var configFile string

	home, err := os.UserHomeDir()
	if err == nil {
        userConf := filepath.Join(home, ".config", "cheatsh")
        if info, err := os.Stat(userConf); err == nil && info.IsDir() {
        	configFile = filepath.Join(userConf, "commands.json")
        } else {
        	configFile = "/etc/cheatsh/commands.json"
        }
    }

	var commands CmdList

	if *configFlag != "" {
		commands, err = loadCommands(*configFlag)
		if err != nil {
			log.Fatalf("Cannot load commands file: %v", err)
		}
		StartTui(commands)
		return
	} else if len(os.Args) == 1 {
		commands, err = loadCommands(configFile)
		if err != nil {
			log.Fatalf("Cannot load commands file: %v", err)
		}
		StartTui(commands)
		return
	}

	switch {

		case *helpFlag:
			printHelp()

		case *newFlag:
			HandleInput()

		default:
			printHelp()
			os.Exit(1)

	}

}

func printHelp() {
	fmt.Println(`Usage: cheatsh [options]
	Options:
	  --config <file>  Specify a config file
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
