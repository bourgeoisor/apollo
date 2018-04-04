package main

import (
	"sort"
	"strings"
)

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

// Removes common words from the prefix of title strings.
func cleanTitlePrefix(str string) string {
	newStr := str

	newStr = strings.Replace(newStr, "The ", "", 1)
	newStr = strings.Replace(newStr, "A ", "", 1)

	return newStr
}

// Sorting function used to sort entries by title.
func titleSortFunc(e1, e2 *Entry) bool {
	s1 := cleanTitlePrefix(e1.Title)
	if e1.TitleSort != "" {
		s1 = cleanTitlePrefix(e1.TitleSort)
	}

	s2 := cleanTitlePrefix(e2.Title)
	if e2.TitleSort != "" {
		s2 = cleanTitlePrefix(e2.TitleSort)
	}

	if strings.ToLower(s1) == strings.ToLower(s2) {
		return e1.Year < e2.Year
	} else {
		return strings.ToLower(s1) < strings.ToLower(s2)
	}
}

// Sorting function used to sort entries by year.
func yearSortFunc(e1, e2 *Entry) bool {
	return e1.Year < e2.Year
}

// Sorting function used to sort entries by rating.
func ratingSortFunc(e1, e2 *Entry) bool {
	return e1.Rating > e2.Rating
}

// Sorting function used to sort entries by info field.
func infoSortFunc(e1, e2 *Entry) bool {
	return e1.Info < e2.Info
}