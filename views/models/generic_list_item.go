package models

// Represents a generic item in a list
// Used as part of the custom ItemDelegate implementation
type genericListItem struct {
	title       string
	desc        string
	filterValue string
}

func (i genericListItem) Title() string       { return i.title }
func (i genericListItem) Description() string { return i.desc }
func (i genericListItem) FilterValue() string { return i.filterValue }
