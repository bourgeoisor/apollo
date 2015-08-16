package main

import (
    "strings"
    "log"
)

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
                log.Print(err)
            }
        } else {
            a.logError("Wrong number of arguments.")
        }
    case "/close":
        err := a.closeCurrentTab()
        if err != nil {
            log.Print(err)
        }
    case "/set":
        if len(args) == 3 {
            a.c.set(args[1], args[2])
        } else {
            a.logError("Wrong number of arguments.")
        }
    case "/config":
        for _, value := range a.c.config() {
            a.log(value)
        }
    default:
        a.logError("'" + command + "' is not a valid command.")
    }
}
