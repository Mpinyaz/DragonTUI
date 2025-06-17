package views

import (
	"DragonTUI/internal/models"
	"DragonTUI/internal/utils"
	"fmt"
	"regexp"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	// "github.com/charmbracelet/lipgloss"
)

type ContactModel struct {
	Form        *huh.Form
	Width       int
	Height      int
	Help        help.Model
	KeyMap      ContactKeyMap
	FeedbackMsg FeedbackMsg
}

type FeedbackMsg struct {
	email   string
	name    string
	message string
}

func newForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("CodeDragon Mailer").
				Description("Full Name: ").
				Placeholder("Enter your full name here").
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Name required, try again")
					}
					return nil
				}).
				Key("name"),
			huh.NewInput().
				Description("Email: ").
				Placeholder("Enter your email here").
				Key("email").
				Validate(func(str string) error {
					re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
					if str != "" {
						if !re.MatchString(str) {
							return fmt.Errorf("invalid Email, try again")
						}
					} else {
						return fmt.Errorf("email required, try again")
					}
					return nil
				}),
			huh.NewText().
				CharLimit(300).
				Key("message").
				Description("Mailbox:").
				Placeholder("Enter message here").
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Message required, try again")
					}
					return nil
				}).
				Lines(5),
			huh.NewConfirm().Key("done").Title("Send Message?").Affirmative("Yes!").Negative("Cancel"),
		),
	)
}

func NewContactModel(width int, height int) *ContactModel {
	contactform := newForm()

	return &ContactModel{Width: width, Height: height, Help: help.New(), KeyMap: ContactKeyMap{}, Form: contactform}
}

type ContactKeyMap struct{}

func (k ContactKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Back to Menu")),
	}
}

func (k ContactKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func (m *ContactModel) updateDimensions(width, height int) {
	m.Width = width
	m.Height = height
}

func (m *ContactModel) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("Contact Me"), m.Form.Init())
}

func (m *ContactModel) Update(msg tea.Msg) (models.Page, tea.Cmd) {
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
			m.Form = newForm()
			var cmd tea.Cmd
			// s := GetMenuModel(m.Width, m.Height)
			// s.Init()
			// return s, tea.Batch(cmd, CheckWeather)
			return GetMenuModel(m.Width, m.Height), tea.Batch(cmd, tea.SetWindowTitle("Dragon's Lair"), CheckWeather)
		case " ":
			return m, cmd

		}

	case tea.WindowSizeMsg:
		m.updateDimensions(msg.Width, msg.Height)
		return m, cmd
	}

	// Process the form
	form, cmd := m.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Form = f
		cmds = append(cmds, cmd)
	}

	if m.Form.State == huh.StateCompleted {
		if m.Form.GetBool("done") {
			m.FeedbackMsg = FeedbackMsg{
				email:   m.Form.GetString("email"),
				name:    m.Form.GetString("name"),
				message: m.Form.GetString("message"),
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *ContactModel) View() string {
	switch m.Form.State {
	case huh.StateCompleted:
		if m.Form.State == huh.StateNormal {
			return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, m.Form.View())
		}
		s := lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render(fmt.Sprintf("\nHey %s, message was delivered \n", utils.Rainbow(lipgloss.NewStyle(), m.FeedbackMsg.name, blends)))
		hlp := fmt.Sprintf("\n%s", m.Help.View(m.KeyMap))

		finalRender := fmt.Sprintf("\n%s\n\n%s", banner.Bold(true).Italic(true).Render(s), hlp)
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, finalRender)

	default:
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, m.Form.View())
	}
}
