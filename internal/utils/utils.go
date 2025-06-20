package utils

import (
	"image/color"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const Logo string = `
  ██╗    ██╗██████╗ ██████╗  █████╗  ██████╗  ██████╗ ███╗   ██╗██╗  
 ██╔╝   ██╔╝██╔══██╗██╔══██╗██╔══██╗██╔════╝ ██╔═══██╗████╗  ██║╚██╗ 
██╔╝   ██╔╝ ██║  ██║██████╔╝███████║██║  ███╗██║   ██║██╔██╗ ██║ ╚██╗
╚██╗  ██╔╝  ██║  ██║██╔══██╗██╔══██║██║   ██║██║   ██║██║╚██╗██║ ██╔╝
 ╚██╗██╔╝   ██████╔╝██║  ██║██║  ██║╚██████╔╝╚██████╔╝██║ ╚████║██╔╝ 
  ╚═╝╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝╚═╝  
	`

func Rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}

func ColorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := range make([]int, ySteps) {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := range make([]int, xSteps) {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
