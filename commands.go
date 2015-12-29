package main

import (
	"strings"
)

// HandleCommand takes the latest user input, parses it, and calls the wanted function.
func (a *Apollo) handleCommand() {
	args := strings.Split(string(a.input), " ")
	command := args[0]

	switch command {
	case "/quit":
		a.running = false
	case "/help":
		if len(args) == 1 {
			a.printHelp()
		} else {
			a.printDetailedHelp(args[1])
		}
	case "/open":
		if len(args) == 2 {
			err := a.openTab(args[1])
			if err != nil {
				a.logError(err.Error())
			}
		} else {
			a.logError("term: invalid number of arguments")
		}
	case "/close":
		err := a.closeCurrentTab()
		if err != nil {
			a.logError(err.Error())
		}
	case "/set":
		if len(args) == 3 {
			err := a.c.set(args[1], args[2])
			if err != nil {
				a.logError(err.Error())
			} else {
				a.log("{b}│ {d}Configuration changed.")
			}
		} else {
			a.logError("term: invalid number of arguments")
		}
	case "/config":
		a.printConfig()
	case "/stats":
		a.printStats()
	default:
		a.logError("term: invalid command")
	}
}
