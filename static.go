package main

// Version is the version number of the application.
const version = "Apollo v.0.3.5"

// PrintHelp prints out the help guide to the logs.
func (a *Apollo) printHelp() {
	s := []string{
		"{b}*───( Main Help Guide )───*",
		"{b}│ {d}List of key-bindings:",
		"{b}│ {d}Ctrl+C.........Close this software",
		"{b}│ {d}Alt+Num........Go to the [num]th tab",
		"{b}│ {d}Enter..........Send, or toggle between input bar and tab",
		"{b}│",
		"{b}│ {d}List of commands:",
		"{b}│ {d}(For more details, use /help <command name>)",
		"{b}│ {d}/help..........Show this help guide",
		"{b}│ {d}/quit..........Close this software",
		"{b}│ {d}/open..........Create a new tab",
		"{b}│ {d}/close.........Close the current tab",
		"{b}│ {d}/set...........Set a configuration option",
		"{b}│ {d}/config........Show the current configuration",
		"{b}*───*",
	}

	for i := 0; i < len(s); i++ {
		a.log(s[i])
	}
}

// PrintDetailedHelp prints out the detailed help of a function to the logs.
func (a *Apollo) printDetailedHelp(subject string) {
	var s []string
	switch subject {
	case "help":
		s = []string{
			"{b}*───( Detailed Help Guide )───*",
			"{b}│ {d}/help",
			"{b}│ {d}/help <command name>",
			"{b}│",
			"{b}│ {d}Displays a help guide.",
			"{b}*───*",
		}
	case "quit":
		s = []string{
			"{b}*───( Detailed Help Guide )───*",
			"{b}│ {d}/quit",
			"{b}│",
			"{b}│ {d}Closes this software.",
			"{b}*───*",
		}
	case "open":
		s = []string{
			"{b}*───( Detailed Help Guide )───*",
			"{b}│ {d}/open <tab name>",
			"{b}│",
			"{b}│ {d}Opens a given tab or, if it already exists, selects it.",
			"{b}│ {d}Tabs available: anime, books, games, movies, series.",
			"{b}│ {d}For keybindings, use '/help tabs'",
			"{b}*───*",
		}
	case "close":
		s = []string{
			"{b}*───( Detailed Help Guide )───*",
			"{b}│ {d}/close",
			"{b}│",
			"{b}│ {d}Close the current tab.",
			"{b}*───*",
		}
	case "set":
		s = []string{
			"{b}*───( Detailed Help Guide )───*",
			"{b}│ {d}/set <option> <value>",
			"{b}│",
			"{b}│ {d}Sets a configuration option to a specific value.",
			"{b}*───*",
		}
	case "config":
		s = []string{
			"{b}*───( Detailed Help Guide )───*",
			"{b}│ {d}/config",
			"{b}│",
			"{b}│ {d}Shows the configuration options and their values.",
			"{b}*───*",
		}
	case "tabs":
		s = []string{
			"{b}*───( Detailed Help Guide )───*",
			"{b}│ {d}List of key-bindings:",
			"{b}│ {d}1..............Switch to passive view",
			"{b}│ {d}2..............Switch to active view",
			"{b}│ {d}3..............Switch to inactive view",
			"{b}│ {d}4..............Switch to all view",
			"{b}│ {d}s..............Sort the entries",
			"{b}│ {d}D..............Delete the current entry",
			"{b}│ {d}e..............Edit the current entry",
			"{b}│ {d}t..............Tag the current entry",
			"{b}│ {d}r..............Toggle ratings",
			"{b}│ {d}a..............Toggle the current entry's state",
			"{b}│ {d}z/x............Change the rating of the current entry",
			"{b}│ {d}c/v............Change the current episode of an entry",
			"{b}│ {d}p..............Print the current entries to a file",
			"{b}*───*",
		}
	default:
		s = []string{
			"{b}│ {d}Detailed help does not exist for this command.",
		}
	}

	for i := 0; i < len(s); i++ {
		a.log(s[i])
	}
}

// PrintWelcome prints out the welcome message to the logs.
func (a *Apollo) printWelcome() {
	a.log("{b}*───( " + version + " )───*")
	a.log("{b}│ {d}This software is under heavy developpment and may contain bugs and glitches.")
	a.log("{b}│ {d}Use at your own risk. To get started, use /help.")
	a.log("{b}*───*")
}

// PrintConfig prints out the list of configuration options to the logs.
func (a *Apollo) printConfig() {
	a.log("{b}*───( Current Configuration )───*")
	for _, value := range a.c.config() {
		a.log("{b}│ {d}" + value)
	}
	a.log("{b}*───*")
}
