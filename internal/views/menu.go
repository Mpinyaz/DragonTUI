package views

import (
	"DragonTUI/internal/models"
	"DragonTUI/internal/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuModel struct {
	Text             string
	Quitting         bool
	AltScreen        bool
	Suspending       bool
	Width            int
	Height           int
	Spinner          spinner.Model
	Help             help.Model
	KeyMap           MenuKeyMap
	SelectedMenuItem string
	MenuList         list.Model
}

type item struct{ title, desc string }

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m *MenuModel) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("Dragon's Lair"), m.Spinner.Tick)
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
		case "enter":
			i, ok := m.MenuList.SelectedItem().(item)
			if ok {
				m.SelectedMenuItem = string(i.title)
			}

			if i.title == "About" {
				var cmd tea.Cmd
				s := GetAboutModel(m.Width, m.Height)
				s.Init()
				return s, cmd

			}
			if i.title == "Contact Me" {
				var cmd tea.Cmd
				s := GetContactModel(m.Width, m.Height)
				s.Init()
				return s, cmd

			}
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.MenuList, cmd = m.MenuList.Update(msg)
	return m, cmd
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

func NewMenuModel(width, height int) *MenuModel {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#edff82"))
	items := []list.Item{
		item{title: "About", desc: "Find out more about my skills and experience"},
		item{title: "Contact Me", desc: "Send me an email!!!"},
		item{title: "Github Repo", desc: "Explore my side projects"},
	}

	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(lipgloss.Color("#fffdf5")).Background(lipgloss.Color("#6f03fc")).Bold(true).Blink(true)

	l := list.New(items, d, 80, 15)
	l.Styles.Title = list.DefaultStyles().Title.Padding(0).Margin(1)
	l.Title = "Learn more about me"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return &MenuModel{
		Text:       "Loading App...press q to quit",
		Quitting:   false,
		AltScreen:  true,
		Suspending: false,
		Width:      width,
		Height:     height,
		Spinner:    s,
		Help:       help.New(),
		KeyMap:     MenuKeyMap{},
		MenuList:   l,
	}
}

func (m MenuModel) View() string {
	// if m.Width == 0 {
	// 	return fmt.Sprintf("\n\n\t%s %s\n\n", m.Spinner.View(), lipgloss.NewStyle().Render(utils.Rainbow(lipgloss.NewStyle(), m.Text, blends)))
	// }
	// if m.Quitting {
	// 	return fmt.Sprintf("Bye \n")
	// }
	m.Help.Styles.ShortDesc = style.Faint(true).Blink(true)
	m.Help.ShortSeparator = " • "
	m.Help.Styles.ShortSeparator = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("#334dcc"))
	m.Help.Styles.ShortKey = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#fff7db"))
	banner := fmt.Sprintf("\n%s\n", banner.Render(utils.Rainbow(lipgloss.NewStyle(), utils.Logo, blends)))
	menuList := m.MenuList.View()
	keymap := fmt.Sprintf("\n\n%s\n", m.Help.View(m.KeyMap))

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, banner+menuList+keymap)
}
