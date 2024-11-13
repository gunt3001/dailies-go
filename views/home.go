package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
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
}

// NewHomeView Constructor
func NewHomeView() HomeView {
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

	// On keypress
	case tea.KeyMsg:
		switch {
		// q to quit
		case key.Matches(msg, v.keymap.quit):
			return v, tea.Quit
		}
	}

	// Return the updated model and no command
	return v, nil
}

func (v HomeView) View() string {

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())
	borderSizeX := borderStyle.GetBorderLeftSize() + borderStyle.GetBorderRightSize()
	borderSizeY := borderStyle.GetBorderTopSize() + borderStyle.GetBorderBottomSize()

	mainAreaStyle := borderStyle.
		Width(v.width - v.sidebarWidth - borderSizeX).
		Height(v.height - borderSizeY)

	sidebarStyle := borderStyle.
		Width(v.sidebarWidth - borderSizeX).
		Height(v.height - borderSizeY)

	helpView := v.help.ShortHelpView([]key.Binding{
		v.keymap.quit,
	})

	ui := lipgloss.JoinHorizontal(lipgloss.Top,
		mainAreaStyle.Render("Main Area"),
		sidebarStyle.Render("Sidebarrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr")) +
		"\n" + helpView

	// Return the UI for rendering
	return ui
}
