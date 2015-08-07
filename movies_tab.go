package main

import (
    "github.com/nsf/termbox-go"
    "strconv"
)

type MoviesTab struct {
    a *Apollo
    name string
    status string

    movies []Movie
    view string
    offset int
    cursor int
}

func CreateMoviesTab(a *Apollo) *MoviesTab {
    t := &MoviesTab{
        a: a,
        name: "movies",
        status: "movies",
        view: "watched",
    }

    t.refreshSlice()

    return t
}

func (t *MoviesTab) Name() string {
    return t.name
}

func (t *MoviesTab) Status() string {
    return t.status
}

func (t *MoviesTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Ch {
    case '1':
        t.view = "watched"
        t.refreshSlice()
    case '2':
        t.view = "unwatched"
        t.refreshSlice()
    case '3':
        t.view = "all"
        t.refreshSlice()
    case '4':
        t.view = "rank"
    }

    switch ev.Key {
    case termbox.KeyArrowUp:
        t.cursor--
        if t.cursor < 0 {
            t.cursor = 0
        } else if t.cursor + t.offset < 0 {
            t.offset++
        }
    case termbox.KeyPgup:
        t.cursor -= 5
        if t.cursor < 0 {
            t.cursor = 0
            t.offset = 0
        } else if t.cursor + t.offset < 0 {
            t.offset += 5
        }
    case termbox.KeyArrowDown:
        t.cursor++
        if t.cursor > len(t.movies) - 1 {
            t.cursor--
        } else if t.cursor + t.offset > t.a.height - 4 {
            t.offset--
        }
    case termbox.KeyPgdn:
        t.cursor += 5
        if t.cursor > len(t.movies) - 1 {
            t.cursor = len(t.movies) - 1
            t.offset = -len(t.movies) + t.a.height - 4
        } else if t.cursor + t.offset > t.a.height - 4 {
            t.offset -= 5
        }
    default:
        return false
    }

    return true
}

func (t *MoviesTab) Draw() {
    for j := 0; j < t.a.height - 2; j++ {
        if j < len(t.movies) {
            runes := []rune(t.movies[j + t.offset].Title)
            for i := 0; i < len(runes); i++ {
                termbox.SetCell(i + 3, j + 1, runes[i], color['d'], color['d'])
            }
        }
    }

    termbox.SetCell(1, t.cursor + t.offset + 1, '*', color['d'], color['d'])
}

func (t *MoviesTab) Query(query string) {
    t.a.d.Movies = append(t.a.d.Movies, Movie{Title: query, State: "Watched"})
    t.a.d.save()
}

func (t *MoviesTab) refreshSlice() {
    t.cursor = 0
    t.offset = 0
    t.movies = t.movies[:0]
    for i := 0; i < len(t.a.d.Movies); i++ {
        if t.a.d.Movies[i].State == "Watched" {
            if t.view == "watched" || t.view == "all" {
                t.movies = append(t.movies, t.a.d.Movies[i])
            }
        } else if t.a.d.Movies[i].State == "Unwatched" {
            if t.view == "unwatched" || t.view == "all" {
                t.movies = append(t.movies, t.a.d.Movies[i])
            }
        }
    }

    t.status = "movies - " + t.view + " (" + strconv.Itoa(len(t.movies)) + " entries)"
}
