package main

import (
	"fmt"
	"os"
	"flag"
)


func main() {

	configFlag := flag.String("config", "", "Specify a config file")
	helpFlag := flag.Bool("help", false, "Show help")

	flag.Parse()

	if len(os.Args) == 1 {
		StartTui()
		return
	}

	switch {
		case *helpFlag:
			printHelp()
		case *configFlag != "":
			StartTui(*configFlag)
		default:
			printHelp()
			os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Usage: cheatsh [options]
	Options:
	  --config <file>  Specify a config file
	  --help           Show this help message`)
}
