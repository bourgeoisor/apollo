package main

import (
    "strings"
)

func (a *Apollo) handleCommand() {
    args := strings.Split(string(a.input), " ")
    switch args[0] {
    case "/help":
        a.tabs[0].Query("Help is not made yet okay!")
    default:
        a.tabs[0].Query("Not a valid command..")
    }
}
