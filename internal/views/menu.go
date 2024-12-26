package views

import (
	"DragonTUI/internal/models"
	"DragonTUI/internal/utils"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/gamut"
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
	KeyMap     KeyMap
}

func (m *MenuModel) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m *MenuModel) Update(msg tea.Msg) (models.Page, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = 0
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
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

type KeyMap struct{}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next")),
		key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}
func (k KeyMap) FullHelp() [][]key.Binding {
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
		KeyMap:     KeyMap{},
	}
}

var (
	normal       = lipgloss.Color("#EEEEEE")
	base         = lipgloss.NewStyle().Foreground(normal)
	customBorder = lipgloss.Border{
		Top:          "▀",
		Bottom:       "▄",
		Left:         "█",
		Right:        "█",
		TopLeft:      "╔",
		TopRight:     "╗",
		BottomLeft:   "╚",
		BottomRight:  "╝",
		MiddleLeft:   "╠",
		MiddleRight:  "╣",
		Middle:       "╬",
		MiddleTop:    "╦",
		MiddleBottom: "╩",
	}
	colors = utils.ColorGrid(1, 5)

	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")
	blends = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 50)
	style  = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#e60080")).
		AlignHorizontal(lipgloss.Center).
		MarginLeft(5).
		Blink(true).
		Border(customBorder).
		BorderForeground(lipgloss.Color("#643aff")).
		Padding(1, 3)
	banner = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		MarginLeft(5).
		Border(customBorder).
		BorderForeground(lipgloss.Color("#643aff")).
		Padding(1, 3)
)

func (m MenuModel) View() string {
	var (
		title strings.Builder
	)
	for i, v := range colors {
		const offset = 2
		c := lipgloss.Color(v[0])
		fmt.Fprint(&title, titleStyle.MarginLeft(i*offset).Background(c))
		if i < len(colors)-1 {
			title.WriteRune('\n')
		}
	}

	if m.Width == 0 {
		return fmt.Sprintf("\n\n\t%s %s\n\n", m.Spinner.View(), lipgloss.NewStyle().Render(utils.Rainbow(lipgloss.NewStyle(), m.Text, blends)))
	}
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
	s := fmt.Sprintf("%s\n\n\n%s\n\n", title.String(), banner.Blink(true).Render(utils.Rainbow(lipgloss.NewStyle(), utils.Logo, blends)))
	s += fmt.Sprintf("\n%s", m.Help.View(m.KeyMap))
	return s
}
