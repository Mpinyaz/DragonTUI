package pages

import (
	"fmt"

	"DragonTUI/internal/utils"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	MenuItem int
)

const (
	MenuItemNone MenuItem = iota
	MenuItemAbout
	MenuItemContact
	MenuItemGithub
)

func (m MenuItem) String() string {
	switch m {
	case MenuItemAbout:
		return "About"
	case MenuItemContact:
		return "Contact Me"
	case MenuItemGithub:
		return "Github Repo"
	default:
		return "None"
	}
}

func MenuItemFromString(s string) MenuItem {
	switch s {
	case "About":
		return MenuItemAbout
	case "Contact Me":
		return MenuItemContact
	case "Github Repo":
		return MenuItemGithub
	default:
		return MenuItemNone
	}
}

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
	SelectedMenuItem MenuItem
	MenuList         list.Model
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m *MenuModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Dragon's Lair")
}

func (m *MenuModel) Update(msg tea.Msg) (Page, tea.Cmd) {
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
				m.SelectedMenuItem = MenuItemFromString(i.title)
			}
			switch m.SelectedMenuItem {
			case MenuItemAbout:
				var cmd tea.Cmd
				s := GetAboutModel(m.Width, m.Height)
				s.Init()
				return s, tea.Batch(cmd, tea.SetWindowTitle("About Me"))
			case MenuItemContact:
				var cmd tea.Cmd
				s := GetContactModel(m.Width, m.Height)
				s.Init()
				return s, tea.Batch(cmd, tea.SetWindowTitle("Contact Me"))
			case MenuItemGithub:
				// TODO: Handle Github repo navigation
				return m, tea.Quit
			default:
				return m, tea.Quit
			}
		}
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
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "move down")),
		key.NewBinding(key.WithKeys("esc", "q", "ctrl+c"), key.WithHelp("esc", "exit")),
	}
}

func (k MenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func NewMenuModel(width, height int) *MenuModel {
	sp := spinner.New()
	sp.Spinner = spinner.Globe
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#edff83"))

	items := []list.Item{
		item{title: "About", desc: "Find out more about my skills and experience"},
		item{title: "Contact Me", desc: "Send me an email!!!"},
		item{title: "Github Repo", desc: "Explore my side projects"},
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#fffdf6")).
		Background(lipgloss.Color("#7f03fc")).
		Bold(true).
		Blink(true)

	menuList := list.New(items, d, 81, 15)
	menuList.Styles.Title = list.DefaultStyles().Title.Margin(1)
	menuList.Title = "Learn more about me"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	menuList.SetShowHelp(false)

	return &MenuModel{
		Text:             "Loading App...press q to quit",
		Quitting:         false,
		AltScreen:        true,
		Suspending:       false,
		Width:            width,
		Height:           height,
		Spinner:          sp,
		Help:             help.New(),
		KeyMap:           MenuKeyMap{},
		MenuList:         menuList,
		SelectedMenuItem: MenuItemNone,
	}
}

func (m MenuModel) View() string {
	m.Help.Styles.ShortDesc = utils.Style.Faint(true).Blink(true)
	m.Help.ShortSeparator = " • "
	m.Help.Styles.ShortSeparator = lipgloss.NewStyle().
		Blink(true).
		Foreground(lipgloss.Color("#335dcc"))
	m.Help.Styles.ShortKey = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#fff8db"))

	banner := fmt.Sprintf("\n%s\n", utils.Banner.Render(utils.Rainbow(lipgloss.NewStyle().Blink(true), utils.Logo, utils.Blends)))
	menuList := m.MenuList.View()
	keymap := fmt.Sprintf("\n%s\n", m.Help.View(m.KeyMap))

	finalRender := banner + menuList + keymap

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		finalRender,
	)
}
