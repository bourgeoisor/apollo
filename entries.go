package main

import (
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// EntriesTab is a tab for displaying entries implementing Tabber.
type EntriesTab struct {
	a               *Apollo
	entries         *[]Entry
	slice           []*Entry
	name            string
	status          string
	sortField       string
	view            string
	pastView        string
	offset          int
	cursor          int
	ratings         bool
	search          []Entry
	additionalField string
	entryType       string
	taggingAPI      string
}

// NewEntriesTab creates a new EntriesTab and returns it.
func newEntriesTab(a *Apollo, entries *[]Entry, name string, entryType string, additionalField string, taggingAPI string) *EntriesTab {
	t := &EntriesTab{
		a:               a,
		entries:         entries,
		name:            name,
		sortField:       "title",
		view:            "passive",
		entryType:       entryType,
		additionalField: additionalField,
		taggingAPI:      taggingAPI,
	}

	if t.a.c.get("rating-startup") == "true" {
		t.ratings = true
	}

	t.refreshSlice()

	return t
}

// Name returns the name of the tab.
func (t *EntriesTab) Name() string {
	return t.name
}

// Status returns the status of the tab.
func (t *EntriesTab) Status() string {
	return t.status
}

// ChangeView changes the current view for a specified one.
func (t *EntriesTab) changeView(view string) {
	t.cursor = 0
	t.offset = 0
	t.view = view
	t.refreshSlice()
}

// ToggleSort switches between the different sorts.
func (t *EntriesTab) toggleSort() {
	if t.sortField == "title" {
		t.sortField = "year"
	} else if t.sortField == "year" {
		t.sortField = "rating"
	} else if t.sortField == "rating" && t.additionalField != "" {
		t.sortField = t.additionalField
	} else {
		t.sortField = "title"
	}

	t.cursor = 0
	t.offset = 0
	t.sort()
}

// PrintEntriesToFile saves the current view of entries to a file.
func (t *EntriesTab) printEntriesToFile() {
	var cont string
	for j := 0; j < len(t.slice); j++ {
		year := t.slice[j].Year
		if year == "" {
			year = "----"
		}
		title := t.slice[j].Title

		var str string
		if t.entryType == "additional" {
			info := t.slice[j].Info1
			str = year + " " + title + " [" + info + "]"
		} else if t.entryType == "episodic" {
			episodeDone := strconv.Itoa(t.slice[j].EpisodeDone)
			if len(episodeDone) == 1 {
				episodeDone = "00" + episodeDone
			} else if len(episodeDone) == 2 {
				episodeDone = "0" + episodeDone
			}
			episodeTotal := strconv.Itoa(t.slice[j].EpisodeTotal)
			if len(episodeTotal) == 1 {
				episodeTotal = "00" + episodeTotal
			} else if len(episodeTotal) == 2 {
				episodeTotal = "0" + episodeTotal
			}
			episodes := "[" + episodeDone + "/" + episodeTotal + "]"
			str = episodes + " " + year + " " + title
		} else if t.entryType == "default" {
			str = year + " " + title
		}

		cont += str + "\n"
	}

	path := os.Getenv("HOME") + "/apollo_print.txt"
	err := ioutil.WriteFile(path, []byte(cont), 0644)
	if err != nil {
		t.a.logError(err.Error())
	}
}

// HandleKeyEvent takes a key input and process it, calling the correct function.
func (t *EntriesTab) HandleKeyEvent(ev *termbox.Event) {
	if t.view == "edit" {
		switch ev.Ch {
		case 'e':
			t.view = t.pastView
			t.a.d.save()
			t.refreshSlice()
		case '0':
			t.a.inputActive = true
			t.a.input = []rune(":t " + t.slice[t.cursor].Title)
			t.a.inputCursor = len(t.a.input)
		case '1':
			t.a.inputActive = true
			t.a.input = []rune(":y " + t.slice[t.cursor].Year)
			t.a.inputCursor = len(t.a.input)
		case '2':
			t.a.inputActive = true
			t.a.input = []rune(":f " + t.slice[t.cursor].Future)
			t.a.inputCursor = len(t.a.input)
		case '3':
			if t.entryType == "additional" {
				t.a.inputActive = true
				t.a.input = []rune(":i " + t.slice[t.cursor].Info1)
				t.a.inputCursor = len(t.a.input)
			} else if t.entryType == "episodic" {
				t.a.inputActive = true
				t.a.input = []rune(":e " + strconv.Itoa(t.slice[t.cursor].EpisodeTotal))
				t.a.inputCursor = len(t.a.input)
			}
		}
	} else if t.view == "tag" {
		if ev.Ch == 'q' {
			t.search = t.search[:0]
			t.view = t.pastView
		}

		indexes := map[rune]int{'0': 0, '1': 1, '2': 2, '3': 3,
			'4': 4, '5': 5, '6': 6, '7': 7,
			'8': 8, '9': 9}
		if i, exist := indexes[ev.Ch]; exist {
			if i < len(t.search) {
				t.view = t.pastView
				*t.slice[t.cursor] = t.search[i]
				t.refreshSlice()

				for j, e := range t.slice {
					if *e == t.search[i] {
						t.cursor = j
						t.offset = 0
						if t.cursor > t.a.height-4 {
							t.offset = t.cursor - (t.a.height - 4)
						}
					}
				}
				t.search = t.search[:0]
				t.a.d.save()
			}
		}

	} else {
		switch ev.Ch {
		case '1':
			t.changeView("passive")
		case '2':
			t.changeView("active")
		case '3':
			t.changeView("inactive")
		case '4':
			t.changeView("all")
		case 's':
			t.toggleSort()
		case 'D':
			if len(t.slice) > 0 {
				for i := range *t.entries {
					if (*t.entries)[i] == *t.slice[t.cursor] {
						t.a.logDebug("deleted id." + strconv.Itoa(i))
						*t.entries = append((*t.entries)[:i], (*t.entries)[i+1:]...)
						break
					}
				}
				t.a.d.save()
				t.refreshSlice()
			}
		case 'e':
			if len(t.slice) > 0 {
				t.pastView = t.view
				t.view = "edit"
			}
		case 't':
			if len(t.slice) > 0 {
				t.fetchTags()
			}
		case 'r':
			t.ratings = !t.ratings
		case 'a':
			if len(t.slice) > 0 {
				if t.slice[t.cursor].State == "passive" {
					t.slice[t.cursor].State = "inactive"
				} else if t.slice[t.cursor].State == "inactive" {
					t.slice[t.cursor].State = "active"
				} else {
					t.slice[t.cursor].State = "passive"
				}
				t.slice[t.cursor].Rating = 0
				t.a.d.save()
			}
		case 'z':
			if len(t.slice) > 0 && t.ratings {
				if t.slice[t.cursor].Rating > 0 {
					t.slice[t.cursor].Rating--
					t.a.d.save()
				}
			}
		case 'x':
			if len(t.slice) > 0 && t.ratings {
				if t.slice[t.cursor].Rating < 6 {
					t.slice[t.cursor].Rating++
					t.a.d.save()
				}
			}
		case 'c':
			if len(t.slice) > 0 && t.entryType == "episodic" {
				t.slice[t.cursor].State = "active"
				if t.slice[t.cursor].EpisodeDone > 0 {
					t.slice[t.cursor].EpisodeDone--
				} else {
					t.slice[t.cursor].State = "inactive"
				}
				t.a.d.save()
			}
		case 'v':
			if len(t.slice) > 0 && t.entryType == "episodic" {
				if t.slice[t.cursor].EpisodeDone < t.slice[t.cursor].EpisodeTotal ||
					t.slice[t.cursor].EpisodeTotal == 0 {
					t.slice[t.cursor].EpisodeDone++
					t.slice[t.cursor].State = "active"
					if t.slice[t.cursor].EpisodeDone == t.slice[t.cursor].EpisodeTotal {
						t.slice[t.cursor].Rating = 0
						t.slice[t.cursor].State = "passive"
					}
				}
				t.a.d.save()
			}
		case 'p':
			t.printEntriesToFile()
		}

		switch ev.Key {
		case termbox.KeyArrowUp:
			t.cursor--
			if t.cursor < 0 {
				t.cursor = 0
			} else if t.cursor-t.offset < 0 {
				t.offset--
			}
		case termbox.KeyPgup:
			t.cursor -= 5
			if t.cursor < 0 {
				t.cursor = 0
				t.offset = 0
			} else if t.cursor-t.offset < 0 {
				t.offset -= 5
				if t.offset < 0 {
					t.offset = 0
				}
			}
		case termbox.KeyArrowDown:
			t.cursor++
			if t.cursor > len(t.slice)-1 {
				t.cursor--
			} else if t.cursor-t.offset > t.a.height-4 {
				t.offset++
			}
		case termbox.KeyPgdn:
			t.cursor += 5
			if t.cursor > len(t.slice)-1 {
				t.cursor = len(t.slice) - 1
				if len(t.slice) > t.a.height-3 {
					t.offset = len(t.slice) - (t.a.height - 3)
				}
			} else if t.cursor-t.offset > t.a.height-4 {
				t.offset += 5
				if t.cursor-t.offset < t.a.height-3 {
					t.offset = t.cursor - (t.a.height - 4)
				}
			}
		}
	}

	t.a.logDebug("cursor: " + strconv.Itoa(t.cursor) + " offset: " + strconv.Itoa(t.offset) +
		" len:" + strconv.Itoa(len(t.slice)) + " height:" + strconv.Itoa(t.a.height))
}

// DrawEditView draws the edit view on the terminal.
func (t *EntriesTab) drawEditView() {
	t.a.drawString(0, 1, "{b}*───( Editing Entry )───")
	t.a.drawString(0, 2, "{b}│ {C}e. {d}Return to the entry list.")
	t.a.drawString(0, 3, "{b}│")
	t.a.drawString(0, 4, "{b}│ {C}0. {d}[{B}Title{d}]    "+t.slice[t.cursor].Title)
	t.a.drawString(0, 5, "{b}│ {C}1. {d}[{B}Year{d}]     "+t.slice[t.cursor].Year)
	t.a.drawString(0, 6, "{b}│ {C}2. {d}[{B}Future{d}]   "+t.slice[t.cursor].Future)
	if t.entryType == "additional" {
		t.a.drawString(0, 7, "{b}│ {C}3. {d}[{B}Info{d}]     "+t.slice[t.cursor].Info1)
		t.a.drawString(0, 8, "{b}*───*")
	} else if t.entryType == "episodic" {
		t.a.drawString(0, 7, "{b}│ {C}3. {d}[{B}Episodes{d}] "+strconv.Itoa(t.slice[t.cursor].EpisodeTotal))
		t.a.drawString(0, 8, "{b}*───*")
	} else {
		t.a.drawString(0, 7, "{b}*───*")
	}
}

// DrawTagView draws the tag view on the terminal.
func (t *EntriesTab) drawTagView() {
	t.a.drawString(0, 1, "{b}*───( Tagging Entry )───")
	t.a.drawString(0, 2, "{b}│ {C}q. {d}Cancel tagging.")
	t.a.drawString(0, 3, "{b}│")
	for j := 0; j < len(t.search); j++ {
			entry := t.search[j]

		year := entry.Year
		if year == "" {
			year = "----"
		}
		str := "{b}│ {C}" + strconv.Itoa(j) + ". {d}[{B}" + year + "{d}] " + entry.Title
		if t.entryType == "additional" {
			str += " [" + entry.Info1 + "]"
		} else if t.entryType == "episodic" {
			str += " [" + strconv.Itoa(entry.EpisodeTotal) + "]"
		}
		t.a.drawString(0, j+4, str)
	}
	t.a.drawString(0, len(t.search)+4, "{b}*───*")
}

// DrawEntries draws all the entries of the view on the terminal.
func (t *EntriesTab) drawEntries() {
	for j := 0; j < t.a.height-3; j++ {
		if j < len(t.slice) {
			entry := t.slice[j+t.offset]

			i := 3
			if t.ratings && t.view != "active" {
				for i := 0; i < entry.Rating; i++ {
					if entry.State == "passive" {
						termbox.SetCell(i+3, j+1, '*', colors['g'], colors['d'])
					} else if entry.State == "inactive" {
						termbox.SetCell(i+3, j+1, '*', colors['b'], colors['d'])
					}
				}
				i = 10
			}

			colorMod := "{d}"
			if entry.Future != "" {
				colorMod = "{D}"
			}

			year := entry.Year
			if year == "" {
				year = "----"
			}
			switch entry.State {
			case "passive":
				year = "{g}" + year + colorMod
			case "active":
				year = "{y}" + year + colorMod
			case "inactive":
				year = "{b}" + year + colorMod
			}
			title := entry.Title

			var str string
			if t.entryType == "additional" {
				info := entry.Info1
				str = year + " " + title + " [{b}" + info + "{d}]"
			} else if t.entryType == "episodic" {
				episodeDone := strconv.Itoa(entry.EpisodeDone)
				if episodeDone == "0" {
					episodeDone = "{k}000"
				} else if len(episodeDone) == 1 {
					episodeDone = "{k}00{b}" + episodeDone
				} else if len(episodeDone) == 2 {
					episodeDone = "{k}0{b}" + episodeDone
				} else if len(episodeDone) == 3 {
					episodeDone = "{b}" + episodeDone
				}
				episodeTotal := strconv.Itoa(entry.EpisodeTotal)
				if episodeTotal == "0" {
					episodeTotal = "{k}???"
				} else if len(episodeTotal) == 1 {
					episodeTotal = "{k}00{b}" + episodeTotal
				} else if len(episodeTotal) == 2 {
					episodeTotal = "{k}0{b}" + episodeTotal
				} else if len(episodeTotal) == 3 {
					episodeTotal = "{b}" + episodeTotal
				}
				episodes := "[" + episodeDone + "{d}/" + episodeTotal + "{d}]"
				str = episodes + " " + year + " " + title
			} else if t.entryType == "default" {
				str = year + " " + title
			}

			t.a.drawString(i, j+1, str)

			if entry.Future != "" {
				futureStr := "{D}(" + entry.Future + ")"
				t.a.drawStringRightAlign(t.a.width, j+1, futureStr)
			}
		}
	}

	termbox.SetCell(1, t.cursor-t.offset+1, '*', colors['d'], colors['d'])
}

// Draw calls the correct drawing function depending on the view.
func (t *EntriesTab) Draw() {
	if t.view == "edit" {
		t.drawEditView()
	} else if t.view == "tag" {
		t.drawTagView()
	} else {
		t.drawEntries()
	}
}

// AppendEntry adds a new entry to the list of entries and moves the cursor to it.
func (t *EntriesTab) appendEntry(e Entry) {
	*t.entries = append(*t.entries, e)
	t.a.d.save()
	t.refreshSlice()

	for i, e2 := range t.slice {
		if *e2 == e {
			t.cursor = i
			if t.cursor > t.a.height-4 {
				t.offset = t.cursor - (t.a.height - 4)
			} else {
				t.offset = 0
			}
		}
	}
}

// EditCurrentEntry changes the value of one of the current entry's data.
func (t *EntriesTab) editCurrentEntry(field rune, value string) {
	switch field {
	case 't':
		t.slice[t.cursor].Title = value
	case 'y':
		t.slice[t.cursor].Year = value
	case 'f':
		t.slice[t.cursor].Future = value
	case 'e':
		episodeTotal, err := strconv.Atoi(value)
		if err == nil {
			t.slice[t.cursor].EpisodeTotal = episodeTotal
		}
	case 'i':
		t.slice[t.cursor].Info1 = value
	}

	t.a.inputActive = false
}

// Query processes the user input and calls the correct function.
func (t *EntriesTab) Query(query string) {
	if query[0] == ':' {
		t.editCurrentEntry(rune(query[1]), query[3:])
	} else if query == "!focused" {
		t.cursor = 0
		t.offset = 0
		t.view = "active"
		t.refreshSlice()
		if len(t.slice) == 0 {
			t.cursor = 0
			t.offset = 0
			t.view = "passive"
			t.refreshSlice()
		}
	} else {
		t.appendEntry(Entry{Title: query, State: t.entryState()})
		if t.a.c.get("auto-tag") == "true" {
			t.a.inputActive = false
			t.fetchTags()
		}
	}
}

// By is the helper type used to sort entries.
type By func(e1, e2 *Entry) bool

// Sort creates an EntrySorter and sorts the entries.
func (by By) Sort(entries []*Entry) {
	es := &entrySorter{
		entries: entries,
		by:      by,
	}
	sort.Sort(es)
}

// EntrySorter is the structure that contains the entries to be sorted and the function
// used to sort them.
type entrySorter struct {
	entries []*Entry
	by      func(e1, e2 *Entry) bool
}

// Len returns the length of a slice of entries.
func (s *entrySorter) Len() int {
	return len(s.entries)
}

// Swap changes the position of two different entries in a slice.
func (s *entrySorter) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

// Less returns if an entry is lesser than another entry.
func (s *entrySorter) Less(i, j int) bool {
	return s.by(s.entries[i], s.entries[j])
}

// Sort takes a slice and sorts it according to the sort option toggled.
func (t *EntriesTab) sort() {

	title := func(e1, e2 *Entry) bool {
		if strings.ToLower(e1.Title) == strings.ToLower(e2.Title) {
			return e1.Year < e2.Year
		} else {
			return strings.ToLower(e1.Title) < strings.ToLower(e2.Title)
		}
	}

	year := func(e1, e2 *Entry) bool {
		return e1.Year < e2.Year
	}

	rating := func(e1, e2 *Entry) bool {
		return e1.Rating > e2.Rating
	}

	info1 := func(e1, e2 *Entry) bool {
		return e1.Info1 < e2.Info1
	}

	By(title).Sort(t.slice)

	if t.sortField == "year" {
		By(year).Sort(t.slice)
	} else if t.sortField == "rating" {
		By(rating).Sort(t.slice)
	} else if t.sortField != "title" {
		By(info1).Sort(t.slice)
	}

	t.status = t.name + " - " + t.view + " (" + strconv.Itoa(len(t.slice)) +
		" entries) sorted by " + t.sortField
}

// RefreshSlice sorts the current slice of entries according to which sort is toggled.
func (t *EntriesTab) refreshSlice() {
	t.slice = t.slice[:0]
	for i := range *t.entries {
		if (*t.entries)[i].State == "passive" {
			if t.view == "passive" || t.view == "all" {
				t.slice = append(t.slice, &(*t.entries)[i])
			}
		} else if (*t.entries)[i].State == "inactive" {
			if t.view == "inactive" || t.view == "all" {
				t.slice = append(t.slice, &(*t.entries)[i])
			}
		} else {
			if t.view == "active" || t.view == "all" {
				t.slice = append(t.slice, &(*t.entries)[i])
			}
		}
	}
	t.sort()

	if t.cursor > len(t.slice)-1 {
		t.cursor--
	}

	if t.offset > 0 {
		if len(t.slice)-t.offset < t.a.height-3 {
			t.offset--
		}
	}
}
