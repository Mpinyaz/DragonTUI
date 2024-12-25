package home

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

type AppModel struct {
	text       string
	quitting   bool
	altscreen  bool
	suspending bool
	width      int
	height     int
	spinner    spinner.Model
	help       help.Model
	keymap     keymap
}

const logo string = `


	██╗    ██╗██████╗ ██████╗  █████╗  ██████╗  ██████╗ ███╗   ██╗██╗
	██╔╝   ██╔╝██╔══██╗██╔══██╗██╔══██╗██╔════╝ ██╔═══██╗████╗  ██║╚██╗
	██╔╝   ██╔╝ ██║  ██║██████╔╝███████║██║  ███╗██║   ██║██╔██╗ ██║ ╚██╗
	╚██╗  ██╔╝  ██║  ██║██╔══██╗██╔══██║██║   ██║██║   ██║██║╚██╗██║ ██╔╝
	╚██╗██╔╝   ██████╔╝██║  ██║██║  ██║╚██████╔╝╚██████╔╝██║ ╚████║██╔╝
	╚═╝╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝╚═╝


	`

func (m AppModel) Init() tea.Cmd {

	return m.spinner.Tick
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
	colors = colorGrid(1, 5)

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

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}
func (m AppModel) View() string {
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

	if m.width == 0 {
		return fmt.Sprintf("\n\n\t%s %s\n\n", m.spinner.View(), lipgloss.NewStyle().Render(rainbow(lipgloss.NewStyle(), m.text, blends)))
	}
	if m.quitting == true {
		return fmt.Sprintf("Bye \n")
	}
	m.help.Styles.ShortDesc = style.Faint(true).UnsetBlink()
	m.help.ShortSeparator = " • "
	m.help.Styles.ShortSeparator = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("#334dcc"))
	m.help.Styles.ShortKey = lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Italic(true).
		Foreground(lipgloss.Color("#FFF7DB"))
	s := fmt.Sprintf("%s\n\n\n%s\n\n", title.String(), banner.Blink(true).Render(rainbow(lipgloss.NewStyle(), logo, blends)))
	s += fmt.Sprintf("\n%s", m.help.View(m.keymap))
	return s
}

type keymap struct{}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next")),
		key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
func (s AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type.String() {
		case "ctrl+c", "q", "esc":
			s.quitting = true
			return s, tea.Quit
		case "ctrl+z":
			return s, tea.Suspend
		case " ":
			var cmd tea.Cmd
			if s.altscreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			s.altscreen = !s.altscreen
			return s, cmd
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd

	}

	return s, nil
}
func InitAppModel() AppModel {

	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#edff82"))
	return AppModel{"Loading App...press q to quit", false, true, false, 0, 0, s, help.New(), keymap{}}
}