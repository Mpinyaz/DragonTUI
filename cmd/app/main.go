package main

import (
	"log"
	"os"
	"strconv"

	"DragonTUI/internal/db"
	"DragonTUI/internal/pages"
	"DragonTUI/internal/server"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/joho/godotenv"
)

type appModel struct {
	term          string
	currentPage   pages.Page
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

	if m.currentPage != page {
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

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	initPage := pages.NewMenuModel(pty.Window.Width, pty.Window.Height)

	app := &appModel{
		term:        pty.Term,
		currentPage: initPage,
	}
	return app, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion(), tea.WithOutput(os.Stderr)}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dburl := os.Getenv("DB_URL")
	appPort, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	appHost := os.Getenv("APP_HOST")

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
	server.InitServer(appHost, appPort, teaHandler)
}
