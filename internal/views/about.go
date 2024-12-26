package views

import (
	"DragonTUI/internal/models"
	"DragonTUI/internal/utils"
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
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
}

func (m *AboutModel) Init() tea.Cmd {

	return nil
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
			return NewMenuModel(), cmd
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
		fileContent, err := os.ReadFile("artichoke.md")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not load file:", err)
			os.Exit(1)
		}

		renderedContent, err := glamour.Render(string(fileContent), "dark")
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
	s := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.Viewport.View(), m.footerView())
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, s)
}
func (m *AboutModel) updateDimensions(width, height int) {
	m.Width = width
	m.Height = height
}
func (m *AboutModel) headerView() string {
	var (
		title strings.Builder
	)
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
	return lipgloss.JoinHorizontal(lipgloss.Center, s, lipgloss.NewStyle().Foreground(lipgloss.Color("#e60000")).Render(line))
}

func (m *AboutModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	line := strings.Repeat("─", utils.Max(0, m.Viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, lipgloss.NewStyle().Foreground(lipgloss.Color("#009900")).Render(line), lipgloss.NewStyle().Foreground(lipgloss.Color("#e6f733")).Bold(true).Render(info))
}

func NewAboutModel(width int, height int) *AboutModel {

	return &AboutModel{Width: width, Height: height}
}
