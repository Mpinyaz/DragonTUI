package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Model struct {
	text string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (s Model) View() string {
	textLen := len(s.text)
	Bar := strings.Repeat("*", textLen+4)
	return fmt.Sprintf(
		"%s\n* %s *\n%s\n\nPress Ctrl+c to exit", Bar, s.text, Bar)

}

func (s Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "ctrl+c":
			return s, tea.Quit

		}
	}
	return s, nil
}

func main() {
	p := tea.NewProgram(Model{"This app is under construction"})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running programs:", err)
	}
}
