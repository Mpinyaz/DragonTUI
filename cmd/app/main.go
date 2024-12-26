package main

import (
	"DragonTUI/internal/db"
	"DragonTUI/internal/models"
	"DragonTUI/internal/views"

	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type appModel struct {
	currentPage models.Page
}

func (m *appModel) Init() tea.Cmd {
	if m.currentPage != nil {
		return m.currentPage.Init()
	}
	return nil
}

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.currentPage == nil {
		return m, tea.Quit
	}

	page, cmd := m.currentPage.Update(msg)
	m.currentPage = page
	return m, cmd
}

func (m *appModel) View() string {
	if m.currentPage != nil {
		return m.currentPage.View()
	}
	return "Goodbye!"
}

func main() {
	db, err := db.InitDatabase("dragontui.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	f, err := tea.LogToFile("debug.log", "debug")
	defer db.Close()
	if err != nil {
		log.Fatalf("err: %v", err)

	}
	defer f.Close()
	initPage := views.NewMenuModel()
	app := &appModel{
		currentPage: initPage,
	}
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running programs:", err)
	}
}
