package main

import (
    "github.com/nsf/termbox-go"
    "strconv"
)

type MoviesTab struct {
    a *Apollo
    name string
    status string

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

    return t
}

func (t *MoviesTab) Name() string {
    return t.name
}

func (t *MoviesTab) Status() string {
    return t.status
}

func (a *MoviesTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Key {
    default:
        return false
    }

    return true
}

func (t *MoviesTab) Draw() {
    var movies []Movie
    for i := 0; i < len(t.a.d.Movies); i++ {
        if t.a.d.Movies[i].State == "Watched" {
            if t.view == "all" || t.view == "watched" {
                movies = append(movies, t.a.d.Movies[i])
            }
        } else if t.a.d.Movies[i].State == "Unwatched" {
            if t.view == "all" || t.view == "unwatched" {
                movies = append(movies, t.a.d.Movies[i])
            }
        }
    }

    t.status = "movies - " + t.view + " (" + strconv.Itoa(len(movies)) + " entries)"

    for j := 0; j < t.a.height - 2; j++ {
        if j < len(t.a.d.Movies) {
            runes := []rune(t.a.d.Movies[j + t.offset].Title)
            for i := 0; i < len(runes); i++ {
                termbox.SetCell(i, j + 1, runes[i], termbox.ColorDefault, termbox.ColorDefault)
            }
        }
    }
}

func (t *MoviesTab) Query(query string) {
    t.a.d.Movies = append(t.a.d.Movies, Movie{Title: query, State: "Watched"})
    t.a.d.save()
}
