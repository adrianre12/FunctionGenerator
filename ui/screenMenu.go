package ui

import (
	"TinyGo/FunctionGenerator/text"
	"fmt"

	"tinygo.org/x/tinyfont/proggy"
)

const (
	menuLines = 2
)

type ScreenMenu struct {
	selectedLine uint16
	Text1        *text.Label
	Text2        *text.Label
}

func (s *ScreenMenu) Setup() {
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	s.selectedLine = 1

	text.NewLabel(lcd, font, 0, 7, "Mode")
	s.Text1 = text.NewLabel(lcd, font, 0, 17, fmt.Sprintf("%s", ""))
	s.Text2 = text.NewLabel(lcd, font, 0, 27, fmt.Sprintf("%s", ""))
	//text.NewLabel(lcd, font, 0, 37, "37th Line")
	//text.NewLabel(lcd, font, 0, 47, "01234567890123")

}

func (s *ScreenMenu) Update() {
	s.Text1.Invert = s.selectedLine == 1
	s.Text2.Invert = s.selectedLine == 2
	s.Text1.Write(fmt.Sprintf("%s", "Manual"))
	s.Text2.Write(fmt.Sprintf("%s", "Sweep"))
}

func (s *ScreenMenu) Push(result bool) {
	println("Selected line", s.selectedLine)
	switch s.selectedLine {
	case 1:
		{
			ChangeScreen(Manual)
		}
	default:
		{
			println("Error no screen to select for", s.selectedLine)
		}

	}
}

func (s *ScreenMenu) Rotate(result bool) {
	println("Rotate", result)
	if result && s.selectedLine > 1 { //up
		s.selectedLine--
	}
	if !result && s.selectedLine < menuLines { //down
		s.selectedLine++
	}
}
