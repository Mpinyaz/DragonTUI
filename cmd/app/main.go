package main

import (
	"DragonTUI/internal/db"
	"DragonTUI/internal/home"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

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
	p := tea.NewProgram(home.InitAppModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running programs:", err)
	}
}
