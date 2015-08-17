package main

const version = "Apollo v.0.3.1"

func (a *Apollo) printHelp() {
    s := []string {
        "***",
        "Apollo Help Guide",
        "",
        "List of key-bindings:",
        "Ctrl+C.........Close this software",
        "Alt+Num........Go to the [num]th tab",
        "Enter..........Send, or toggle between input bar and tab",
        "",
        "List of commands:",
        "(For more details, use /help <command name>)",
        "/help..........Show this help guide",
        "/quit..........Close this software",
        "/open..........Create a new tab",
        "/close.........Close the current tab",
        "/set...........Set a configuration option",
        "/config........Show the current configuration",
        "***",
    }

    for i := 0; i < len(s); i++ {
        a.log(s[i])
    }
}

func (a *Apollo) printDetailedHelp(subject string) {
    var s []string
    switch subject {
    case "help":
        s = []string{
            "***",
            "/help",
            "/help <command name>",
            "",
            "Honestly, if you're asking for help about the help command, you should seek for help.",
            "",
            "e.g.",
            "/help",
            "/help quit",
        }
    case "quit":
        s = []string{
            "***",
            "/quit",
            "",
            "Closes this software.",
            "",
            "e.g.",
            "/quit",
        }
    case "open":
        s = []string{
            "***",
            "/open <tab name>",
            "",
            "Opens a given tab or, if it already exists, selects it.",
            "",
            "e.g.",
            "/open movies",
        }
    case "close":
        s = []string{
            "***",
            "/close",
            "",
            "Close the current tab.",
            "",
            "e.g.",
            "/close",
        }
    case "set":
        s = []string{
            "***",
            "/set <option> <value>",
            "",
            "Sets a configuration option to a specific value.",
            "",
            "e.g.",
            "/set autotag false",
        }
    case "config":
        s = []string{
            "***",
            "/config",
            "",
            "Shows the configuration options and their values.",
            "",
            "e.g.",
            "/config",
        }
    default:
        s = []string{
            "***",
            "Detailed help does not exist for this command.",
        }
    }
    s = append(s, "***")

    for i := 0; i < len(s); i++ {
        a.log(s[i])
    }
}

func (a *Apollo) printWelcome() {
    a.tabs[0].Query("{b}*** " + version + " ***")
    a.tabs[0].Query("This software is under heavy developpment and may contain bugs and glitches.")
    a.tabs[0].Query("Use at your own risk. To get started, use /help.")
}
