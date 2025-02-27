package main

import (
	"DragonTUI/internal/db"
	"DragonTUI/internal/models"
	"DragonTUI/internal/views"

	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/joho/godotenv/autoload"
)

var (
	dburl = os.Getenv("DB_URL")
)

type appModel struct {
	currentPage   models.Page
	lastWindowMsg tea.WindowSizeMsg
}

func (m *appModel) Init() tea.Cmd {
	if m.currentPage != nil {
		return m.currentPage.Init()
	}
	return nil
}

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.lastWindowMsg = msg
	}

	if m.currentPage == nil {
		return m, tea.Quit
	}

	page, cmd := m.currentPage.Update(msg)
	m.currentPage = page

	if newPage, ok := m.currentPage.(models.Page); ok && newPage != page {
		return m, tea.Batch(cmd, func() tea.Msg { return m.lastWindowMsg })
	}
	return m, cmd
}

func (m *appModel) View() string {
	if m.currentPage != nil {
		return m.currentPage.View()
	}
	return "Goodbye!"
}

func main() {
	db, err := db.InitDatabase(dburl)
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
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running programs:", err)
	}
}
