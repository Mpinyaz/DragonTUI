package views

import (
	"DragonTUI/internal/models"
	"DragonTUI/internal/utils"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"regexp"
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

func NewContactModel(width int, height int) *ContactModel {
	contactform := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("CodeDragon Mailer \n").
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
							return fmt.Errorf("Invalid Email, try again")
						}
					} else {
						return fmt.Errorf("Email Required, try again")
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
			huh.NewConfirm().
				Key("done").
				Title("Send Message?").
				Affirmative("Yes!").
				Negative("Cancel"),
		),
	)

	return &ContactModel{Width: width, Height: height, Help: help.New(), KeyMap: ContactKeyMap{}, Form: contactform}
}

type ContactKeyMap struct{}

func (k ContactKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "move up")),
		key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "move down")),
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
	return m.Form.Init()
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
			var cmd tea.Cmd
			return GetMenuModel(), cmd
		case " ":
			return m, cmd

		}

	case tea.WindowSizeMsg:
		m.updateDimensions(msg.Width, msg.Height)

	}

	// Process the form
	form, cmd := m.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Form = f
		cmds = append(cmds, cmd)
	}

	if m.Form.State == huh.StateCompleted {
		m.FeedbackMsg = FeedbackMsg{email: m.Form.GetString("email"), name: m.Form.GetString("name"), message: m.Form.GetString("message")}
		// Quit when the form is done.
		// cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)

}
func (m *ContactModel) View() string {
	switch m.Form.State {
	case huh.StateCompleted:
		s := fmt.Sprintf("\nHey %s, message was delivered \n", utils.Rainbow(lipgloss.NewStyle(), m.FeedbackMsg.name, blends))
		return lipgloss.Place(20, 20, lipgloss.Center, lipgloss.Center, banner.Bold(true).Italic(true).Render(s))
	default:
		return m.Form.View()
	}
}
