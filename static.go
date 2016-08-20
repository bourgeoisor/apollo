package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
)

// Version is the version number of the application.
const version = "Apollo v1.0.0"

// Create a configuration directory for Unis systems.
func makeUnixConfigDir() {
	path := os.Getenv("HOME") + "/.config/apollo"
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Print(err)
		}
	}
}

// Returns the path of the database file.
func databasePath() string {
	if runtime.GOOS == "windows" {
		return "database.json"
	} else {
		makeUnixConfigDir()
		return os.Getenv("HOME") + "/.config/apollo/database.json"
	}
}

// Returns the path of the configuration file.
func configurationPath() string {
	if runtime.GOOS == "windows" {
		return "configuration.json"
	} else {
		makeUnixConfigDir()
		return os.Getenv("HOME") + "/.config/apollo/configuration.json"
	}
}

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
		"{b}│ {d}/stats.........Prints some stats about the entries",
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
	case "stats":
		s = []string{
			"{b}*───( Detailed Help Guide )───*",
			"{b}│ {d}/stats",
			"{b}│",
			"{b}│ {d}Prints statistics about the database entries.",
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
	a.log("{b}│ {d}This software is licensed under the MIT License.")
	a.log("{b}│ {d}To get started, use /help.")
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

// PrintStats prints out relevant statistics about the database.
func (a *Apollo) printStats() {
	totalMovies := 0
	watchedMovies := 0
	for i := range a.d.Movies {
		totalMovies++
		if a.d.Movies[i].State == "passive" {
			watchedMovies++
		}
	}

	totalAnime := 0
	watchedAnime := 0
	totalAnimeEp := 0
	watchedAnimeEp := 0
	for i := range a.d.Anime {
		totalAnime++
		if a.d.Anime[i].State == "passive" {
			watchedAnime++
		}
		totalAnimeEp += a.d.Anime[i].EpisodeTotal
		watchedAnimeEp += a.d.Anime[i].EpisodeDone
	}

	totalGames := 0
	playedGames := 0
	for i := range a.d.Games {
		totalGames++
		if a.d.Games[i].State == "passive" {
			playedGames++
		}
	}

	totalBooks := 0
	readBooks := 0
	for i := range a.d.Books {
		totalBooks++
		if a.d.Books[i].State == "passive" {
			readBooks++
		}
	}

	sTotalMovies := strconv.Itoa(totalMovies)
	sWatchedMovies := strconv.Itoa(watchedMovies)
	sTotalAnime := strconv.Itoa(totalAnime)
	sWatchedAnime := strconv.Itoa(watchedAnime)
	sTotalAnimeEp := strconv.Itoa(totalAnimeEp)
	sWatchedAnimeEp := strconv.Itoa(watchedAnimeEp)
	sTotalGames := strconv.Itoa(totalGames)
	sPlayedGames := strconv.Itoa(playedGames)
	sTotalBooks := strconv.Itoa(totalBooks)
	sReadBooks := strconv.Itoa(readBooks)

	a.log("{b}*───( Statistics )───*")
	a.log("{b}│ {d}Movies watched: " + sWatchedMovies + "/" + sTotalMovies)
	a.log("{b}│ {d}Anime seasons watched: " + sWatchedAnime + "/" + sTotalAnime)
	a.log("{b}│ {d} Episodes watched: " + sWatchedAnimeEp + "/" + sTotalAnimeEp)
	a.log("{b}│ {d}Games played: " + sPlayedGames + "/" + sTotalGames)
	a.log("{b}│ {d}Books read: " + sReadBooks + "/" + sTotalBooks)
	a.log("{b}*───*")
}
