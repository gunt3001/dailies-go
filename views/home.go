package views

import (
	"context"
	"dailies-go/db"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	sidebarProportion      = 20  // Width of the sidebar as an integer percent
	helpHeight             = 1   // Height of the help view (in lines)
	maxExpectedEntryLength = 180 // The expected max length of an entry
	maxContentLineCount    = 5   // The maximum number of lines to display content
)

// Struct to handle key bindings
type homeViewKeymap struct {
	quit key.Binding
}

// List item delgate
type itemDelegate struct {
	contentWidth     int
	contentLineCount int
	// Internally rely on existing default delegate
	defaultDelgate list.ItemDelegate
}
type genericListItem struct {
	title       string
	desc        string
	filterValue string
}

func (i genericListItem) Title() string       { return i.title }
func (i genericListItem) Description() string { return i.desc }
func (i genericListItem) FilterValue() string { return i.filterValue }

func (d itemDelegate) Height() int {
	return d.defaultDelgate.Height()
}
func (d itemDelegate) Spacing() int {
	return d.defaultDelgate.Spacing()
}
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.defaultDelgate.Update(msg, m)
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
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

func newItemDelegate(viewWidth int) itemDelegate {
	// Update item height according to view width
	// list view has 2-wide padding on left side
	contentWidth := viewWidth - 2
	if contentWidth < 0 {
		contentWidth = 0
	}
	contentLineCount := getContentLineCount(contentWidth)
	defaultDelegate := list.NewDefaultDelegate()
	defaultDelegate.SetHeight(contentLineCount + 1) // 1 Title line + Content lines

	return itemDelegate{
		contentWidth:     contentWidth,
		defaultDelgate:   defaultDelegate,
		contentLineCount: contentLineCount,
	}
}

// HomeView - Primary view for the application
type HomeView struct {
	width, height, sidebarWidth int
	help                        help.Model
	keymap                      homeViewKeymap
	database                    *db.Queries
	dbContext                   context.Context

	// Styles
	borderStyle   lipgloss.Style
	mainAreaStyle lipgloss.Style
	sidebarStyle  lipgloss.Style

	// Data
	List    list.Model
	entries []db.Entry
}

// NewHomeView Constructor
func NewHomeView(database *db.Queries, context context.Context) HomeView {

	// Initialize entries from db
	entries, err := database.GetEntriesByDateRange(context, db.GetEntriesByDateRangeParams{
		DateStart: "2024-01-01",
		DateEnd:   "2024-12-31",
	})
	// TODO: Handle database read errors
	_ = err
	// Convert into list format
	entryListItems := make([]list.Item, len(entries))
	for i := range entries {
		entryListItems[i] = entries[i]
	}
	// We use default item delegate (item rendering style) here
	// but we customize the height depending on the view width during update cycle
	listView := list.New(entryListItems, newItemDelegate(0), 0, 0)
	listView.SetShowHelp(false)
	listView.SetShowTitle(false)
	listView.SetFilteringEnabled(false)
	listView.SetStatusBarItemName("entry", "entries")
	listView.DisableQuitKeybindings()

	defaultBorderStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	return HomeView{
		width:        0,
		height:       0,
		sidebarWidth: 0,
		help:         help.New(),
		keymap: homeViewKeymap{
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "Quit"),
			),
		},
		database:  database,
		dbContext: context,

		borderStyle:   defaultBorderStyle,
		mainAreaStyle: defaultBorderStyle,
		sidebarStyle:  defaultBorderStyle,

		entries: entries,
		List:    listView,
	}
}

func (v HomeView) Init() tea.Cmd {
	return nil
}

func (v HomeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Determine the type of update message from Bubbletea using switch case
	switch msg := msg.(type) {

	// Note that in Go, cases are not indented within the switch block

	// On Window Size update
	case tea.WindowSizeMsg:
		// Expand the home view to fit the entire window
		v.width = msg.Width
		v.height = msg.Height - helpHeight
		// Recalculate sidebar width
		v.sidebarWidth = v.width * sidebarProportion / 100
		// Update styles
		borderSizeX := v.borderStyle.GetBorderLeftSize() + v.borderStyle.GetBorderRightSize()
		borderSizeY := v.borderStyle.GetBorderTopSize() + v.borderStyle.GetBorderBottomSize()
		mainAreaWidth := v.width - v.sidebarWidth - borderSizeX
		v.mainAreaStyle = v.borderStyle.
			Width(mainAreaWidth).
			Height(v.height - borderSizeY)
		v.sidebarStyle = v.borderStyle.
			Width(v.sidebarWidth - borderSizeX).
			Height(v.height - borderSizeY)
		// Update list size (fill main area)
		v.List.SetSize(v.mainAreaStyle.GetWidth(), v.mainAreaStyle.GetHeight())
		// Update list item size
		v.List.SetDelegate(newItemDelegate(mainAreaWidth))

	// On keypress
	case tea.KeyMsg:
		switch {
		// q to quit
		case key.Matches(msg, v.keymap.quit):
			return v, tea.Quit
		}
	}

	// Update child list view
	var listCmd tea.Cmd
	v.List, listCmd = v.List.Update(msg)

	// Return the updated model and no command
	return v, listCmd
}

func (v HomeView) View() string {

	helpView := v.help.ShortHelpView(append(v.List.ShortHelp(), v.keymap.quit))

	mainAreaView := lipgloss.JoinHorizontal(lipgloss.Top,
		v.mainAreaStyle.Render(v.List.View()),
		v.sidebarStyle.Render("Sidebar content"))

	ui := lipgloss.JoinVertical(lipgloss.Left,
		mainAreaView,
		helpView,
	)

	// Return the UI for rendering
	return ui
}

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
