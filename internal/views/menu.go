package views

import (
	"DragonTUI/internal/models"
	"DragonTUI/internal/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuModel struct {
	Text       string
	Quitting   bool
	AltScreen  bool
	Suspending bool
	Width      int
	Height     int
	Spinner    spinner.Model
	Help       help.Model
	KeyMap     MenuKeyMap
}

func (m *MenuModel) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m *MenuModel) Update(msg tea.Msg) (models.Page, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateDimensions(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		case "ctrl+b":
			var cmd tea.Cmd
			s := GetContactModel(m.Width, m.Height)
			s.Init()
			return s, cmd
		case "ctrl+a":
			var cmd tea.Cmd
			s := GetAboutModel(m.Width, m.Height)
			s.Init()
			return s, cmd
		case " ":
			var cmd tea.Cmd
			if m.AltScreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.AltScreen = !m.AltScreen
			return m, cmd
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}
func (m *MenuModel) updateDimensions(width, height int) {
	m.Width = width
	m.Height = height
}

type MenuKeyMap struct{}

func (k MenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "move up")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "tab enter")),
		key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "move down")),
		key.NewBinding(key.WithKeys("esc", "q", "ctrl+c"), key.WithHelp("esc", "Exit")),
	}
}
func (k MenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func NewMenuModel() *MenuModel {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#edff82"))

	return &MenuModel{
		Text:       "Loading App...press q to quit",
		Quitting:   false,
		AltScreen:  true,
		Suspending: false,
		Width:      0,
		Height:     0,
		Spinner:    s,
		Help:       help.New(),
		KeyMap:     MenuKeyMap{},
	}
}

func (m MenuModel) View() string {
	// if m.Width == 0 {
	// 	return fmt.Sprintf("\n\n\t%s %s\n\n", m.Spinner.View(), lipgloss.NewStyle().Render(utils.Rainbow(lipgloss.NewStyle(), m.Text, blends)))
	// }
	if m.Quitting == true {
		return fmt.Sprintf("Bye \n")
	}
	m.Help.Styles.ShortDesc = style.Faint(true).UnsetBlink()
	m.Help.ShortSeparator = " • "
	m.Help.Styles.ShortSeparator = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("#334dcc"))
	m.Help.Styles.ShortKey = lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Italic(true).
		Foreground(lipgloss.Color("#FFF7DB"))
	s := fmt.Sprintf("\n%s\n\n", banner.Blink(true).Render(utils.Rainbow(lipgloss.NewStyle(), utils.Logo, blends)))
	s += fmt.Sprintf("\n%s", m.Help.View(m.KeyMap))
	return lipgloss.Place(40, 40, lipgloss.Center, lipgloss.Center, s)
}
