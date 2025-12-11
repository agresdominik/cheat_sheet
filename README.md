# Cheat Sheet

A small TUI Tool which displays a pre-configured json file content in a menu and lets the user search for commands and select them. Once selected these are printed in the terminal for easy access.

I wrote this to simplify my search for commands I use often enough to need them multiple times but not often enoigh to remember alway exactly how they are called. 

Written in go + used bubbletea for the TUI part of the application.

## Usage

Use the commands in the given makefile to run, build, test in a docker container or install the application on your device. Beforehand define your own commands.json based on the `commands_template.json` in the `data/` folder.

Once installed just run `cheatsh` in the terminal. If it's not found, make sure `~/.local/bin/` is in your `PATH` or change the path under the `PREFIX` variable in the `Makefile`.


