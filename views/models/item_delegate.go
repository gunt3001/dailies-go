package models

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxExpectedEntryLength = 180 // The expected max length of an entry
	maxContentLineCount    = 5   // The maximum number of lines to display content
)

// List item delgate
// We use this to customize the rendering behavior of the default list item delegate
// where the description line height varies by the width of the list view
type ItemDelegate struct {
	contentWidth     int
	contentLineCount int
	// Internally rely on existing default delegate
	defaultDelgate list.ItemDelegate
}

func (d ItemDelegate) Height() int {
	return d.defaultDelgate.Height()
}
func (d ItemDelegate) Spacing() int {
	return d.defaultDelgate.Spacing()
}
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.defaultDelgate.Update(msg, m)
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	// Automatically wrap text according to available width
	// before passing it to defaultDelegate for normal rendering
	var title, desc, filterValue string
	if i, ok := listItem.(list.DefaultItem); ok {
		title = i.Title()
		desc = i.Description()
		filterValue = i.FilterValue()
	} else {
		return
	}
	wrappedDesc := wrapText(desc, d.contentWidth, d.contentLineCount)
	wrappedListItem := genericListItem{
		title:       title,
		desc:        wrappedDesc,
		filterValue: filterValue,
	}

	d.defaultDelgate.Render(w, m, index, wrappedListItem)
}

// Constructor
func NewItemDelegate(viewWidth int) ItemDelegate {
	// Update item height according to view width
	// list view has 2-wide padding on left side
	contentWidth := viewWidth - 2
	if contentWidth < 0 {
		contentWidth = 0
	}
	contentLineCount := getContentLineCount(contentWidth)
	defaultDelegate := list.NewDefaultDelegate()
	defaultDelegate.SetHeight(contentLineCount + 1) // 1 Title line + Content lines

	return ItemDelegate{
		contentWidth:     contentWidth,
		defaultDelgate:   defaultDelegate,
		contentLineCount: contentLineCount,
	}
}

// Utility Functions

// Calculate the estimated number of lines required to display an entry
// depending on the view width and maximum expected entry content length.
func getContentLineCount(width int) int {
	// Divide the content length by available width to get number of lines
	// but we pad extra spaces (arbitarily chosen to be 5) for when we do
	// word breaking
	availableWidth := width - 5

	// Short-circuit for not enough minimum width
	if availableWidth <= 0 {
		return 1
	}

	lines := (maxExpectedEntryLength / availableWidth) + 1

	if lines > maxContentLineCount {
		return maxContentLineCount
	}
	return lines
}

// Wrap specified string given the available width
// With extra newlines added until minLines is satisfied
func wrapText(text string, width int, minLines int) string {
	if width <= 0 {
		return text // Return the original string if interval is invalid
	}

	wrappedStyle := lipgloss.NewStyle().Width(width).Height(minLines)
	return wrappedStyle.Render(text)
}
