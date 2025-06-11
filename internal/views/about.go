package views

import (
	"DragonTUI/internal/models"
	"DragonTUI/internal/utils"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	aboutTitlestyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return aboutTitlestyle.BorderStyle(b)
	}()
)

type AboutModel struct {
	Content  string
	Ready    bool
	Viewport viewport.Model
	Width    int
	Height   int
	Help     help.Model
	KeyMap   AbtKeyMap
}

func (m *AboutModel) Init() tea.Cmd {
	return tea.SetWindowTitle("About me")
}

func (m *AboutModel) Update(msg tea.Msg) (models.Page, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		case "esc":
			var cmd tea.Cmd
			return GetMenuModel(m.Width, m.Height), cmd
		case " ":
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.updateDimensions(msg.Width, m.Height)
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		if m.Width == 0 {
			m.Viewport = viewport.New(m.Width, m.Height-verticalMarginHeight)
			m.Viewport.YPosition = headerHeight
			m.Viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.Viewport.SetContent(m.Content)
			m.Ready = true

			m.Viewport.YPosition = headerHeight + 1
		} else {
			m.Viewport.Width = m.Width
			m.Viewport.Height = m.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.Viewport))
		}
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *AboutModel) View() string {
	if len(m.Content) == 0 {
		fileContent, err := os.ReadFile("resume.md")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not load file:", err)
			os.Exit(1)
		}

		renderedContent, err := glamour.Render(string(fileContent), "dracula")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not render content:", err)
			os.Exit(1)
		}

		m.Content = renderedContent
	}

	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	if !m.Ready {
		m.Viewport = viewport.New(m.Width, m.Height-verticalMarginHeight)
		m.Viewport.YPosition = headerHeight
		m.Viewport.HighPerformanceRendering = useHighPerformanceRenderer
		m.Viewport.SetContent(m.Content)
		m.Ready = true

		m.Viewport.YPosition = headerHeight + 1
	} else {
		m.Viewport.Width = m.Width
		m.Viewport.Height = m.Height - verticalMarginHeight
	}
	m.Help.Styles.ShortDesc = style.Faint(true).UnsetBlink()
	m.Help.ShortSeparator = " • "
	m.Help.Styles.ShortSeparator = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("#334dcc"))
	m.Help.Styles.ShortKey = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#FFF7DB"))
	s := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.Viewport.View(), m.footerView())
	s += fmt.Sprintf("\n%s", m.Help.View(m.KeyMap))

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, s)
}

func (m *AboutModel) updateDimensions(width, height int) {
	m.Width = width
	m.Height = height
}

func (m *AboutModel) headerView() string {
	var title strings.Builder
	for i, v := range colors {
		const offset = 2
		c := lipgloss.Color(v[0])
		fmt.Fprint(&title, aboutTitlestyle.MarginLeft(i*offset).Background(c))
		if i < len(colors)-1 {
			title.WriteRune('\n')
		}
	}
	s := aboutTitlestyle.Render(utils.Rainbow(lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#ffffff")), "Elton Mpinyuri", blends))
	line := strings.Repeat("─", utils.Max(0, m.Viewport.Width-lipgloss.Width(s)))
	scr := lipgloss.JoinHorizontal(lipgloss.Center, s, lipgloss.NewStyle().Foreground(lipgloss.Color("#e60000")).Render(line))
	return scr
}

func (m *AboutModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	line := strings.Repeat("─", utils.Max(0, m.Viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, lipgloss.NewStyle().Foreground(lipgloss.Color("#009900")).Render(line), lipgloss.NewStyle().Foreground(lipgloss.Color("#e6f733")).Bold(true).Render(info))
}

func NewAboutModel(width int, height int) *AboutModel {
	return &AboutModel{Width: width, Height: height, Help: help.New(), KeyMap: AbtKeyMap{}}
}

type AbtKeyMap struct{}

func (k AbtKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "move up")),
		key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "move down")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Back to Menu")),
	}
}

func (k AbtKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
