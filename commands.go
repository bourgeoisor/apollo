package main

import (
    "strings"
)

func (a *Apollo) handleCommand() {
    args := strings.Split(string(a.input), " ")
    command := args[0]

    switch command {
    case "/quit":
        a.running = false
    case "/help":
        a.log("Help is not made yet okay!")
    case "/open":
        if len(args) == 2 {
            tab := args[1]
            a.openTab(tab)
        } else {
            a.logError("Wrong number of arguments.")
        }
    case "/set":
        if len(args) == 3 {
            option := args[1]
            value := args[2]
            a.c.set(option, value)
        } else {
            a.logError("Wrong number of arguments.")
        }
    case "/close":
        a.closeCurrentTab()
    default:
        a.logError("'" + command + "' is not a valid command.")
    }
}
