package views

import (
	"context"
	"dailies-go/db"
	"dailies-go/views/models"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	sidebarProportion = 20 // Width of the sidebar as an integer percent
	helpHeight        = 1  // Height of the help view (in lines)
)

// Struct to handle key bindings
type homeViewKeymap struct {
	quit key.Binding
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
	listView := list.New(entryListItems, models.NewItemDelegate(0), 0, 0)
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
		v.List.SetDelegate(models.NewItemDelegate(mainAreaWidth))

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
