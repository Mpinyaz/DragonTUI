package views

import (
	"DragonTUI/internal/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/gamut"
)

var (
	normal       = lipgloss.Color("#EEEEEE")
	base         = lipgloss.NewStyle().Foreground(normal)
	customBorder = lipgloss.Border{
		Top:          "▀",
		Bottom:       "▄",
		Left:         "█",
		Right:        "█",
		TopLeft:      "╔",
		TopRight:     "╗",
		BottomLeft:   "╚",
		BottomRight:  "╝",
		MiddleLeft:   "╠",
		MiddleRight:  "╣",
		Middle:       "╬",
		MiddleTop:    "╦",
		MiddleBottom: "╩",
	}
	// colors = [][]string{{"#FF5733"}, {"#33FF57"}, {"#5733FF"}, {"#FFD700"}}

	colors = utils.ColorGrid(1, 5)
	//
	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB"))
	blends = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 50)
	style  = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#e60080")).
		AlignHorizontal(lipgloss.Center).
		MarginLeft(5).
		Blink(true).
		Border(customBorder).
		BorderForeground(lipgloss.Color("#643aff")).
		Padding(1, 3)
	banner = lipgloss.NewStyle().
		Padding(1)
)
