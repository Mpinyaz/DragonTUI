package views

import (
	"fmt"
	"os"
	"regexp"

	"DragonTUI/internal/models"
	"DragonTUI/internal/utils"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/resend/resend-go/v3"
)

type ContactModel struct {
	Form        *huh.Form
	Width       int
	Height      int
	Help        help.Model
	KeyMap      ContactKeyMap
	FeedbackMsg FeedbackMsg
	EmailSent   bool
	EmailError  error
}

type FeedbackMsg struct {
	email   string
	name    string
	message string
}

type EmailSentMsg struct {
	success bool
	err     error
}

func sendEmail(name, email, message string) tea.Cmd {
	return func() tea.Msg {
		resendkey := os.Getenv("RESEND_API_KEY")
		if resendkey == "" {
			return EmailSentMsg{success: false, err: fmt.Errorf("RESEND_API_KEY not set")}
		}

		client := resend.NewClient(resendkey)
		params := &resend.SendEmailRequest{
			From:    "DragonTUI <noreply@resend.dev>",
			To:      []string{"ebmpinyuri@gmail.com"},
			Html:    fmt.Sprintf("<p><strong>Email:</strong> %s</p><p><strong>Message:</strong></p><p>%s</p>", email, message),
			Subject: fmt.Sprintf("New Contact from %s", name),
			ReplyTo: email,
		}

		sent, err := client.Emails.Send(params)
		if err != nil {
			return EmailSentMsg{success: false, err: err}
		}

		fmt.Printf("Email sent successfully! ID: %s\n", sent.Id)
		return EmailSentMsg{success: true, err: nil}
	}
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
						return fmt.Errorf("name required, try again")
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
			huh.NewConfirm().
				Key("done").
				Title("Send Message?").
				Affirmative("Yes!").
				Negative("Cancel"),
		),
	)
}

func NewContactModel(width int, height int) *ContactModel {
	contactform := newForm()
	return &ContactModel{
		Width:      width,
		Height:     height,
		Help:       help.New(),
		KeyMap:     ContactKeyMap{},
		Form:       contactform,
		EmailSent:  false,
		EmailError: nil,
	}
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
	case EmailSentMsg:
		m.EmailSent = msg.success
		m.EmailError = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		case "esc":
			m.Form = newForm()
			m.EmailSent = false
			m.EmailError = nil
			return GetMenuModel(m.Width, m.Height), tea.Batch(cmd, tea.SetWindowTitle("Dragon's Lair"), CheckWeather)
		case " ":
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.updateDimensions(msg.Width, msg.Height)
		return m, cmd
	}

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
			// Send the email
			cmds = append(cmds, sendEmail(m.FeedbackMsg.name, m.FeedbackMsg.email, m.FeedbackMsg.message))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *ContactModel) View() string {
	switch m.Form.State {
	case huh.StateCompleted:
		var s string

		// Show error if email failed
		if m.EmailError != nil {
			s = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Align(lipgloss.Center, lipgloss.Center).
				Render(fmt.Sprintf("\n Error sending email: %v\n\nPlease try again later.\n", m.EmailError))
		} else if m.EmailSent {
			// Show success message
			s = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Center).
				Render(fmt.Sprintf("\n Hey %s, your message was delivered successfully!\n",
					utils.Rainbow(lipgloss.NewStyle(), m.FeedbackMsg.name, blends)))
		} else {
			// Sending in progress
			s = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Center).
				Render("\nSending email...\n")
		}

		hlp := fmt.Sprintf("\n%s", m.Help.View(m.KeyMap))
		finalRender := fmt.Sprintf("\n%s\n\n%s", banner.Bold(true).Italic(true).Render(s), hlp)
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, finalRender)

	default:
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, m.Form.View())
	}
}
