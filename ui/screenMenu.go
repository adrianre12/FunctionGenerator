package ui

import (
	text "TinyGo/FunctionGenerator/Text"
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

	label1 := text.NewLabel(lcd, font, 0, 7, "Line1 ")
	_, labelW := label1.LineWidth()
	s.Text1 = text.NewLabel(lcd, font, int16(labelW), 7, fmt.Sprintf("%s", ""))

	label2 := text.NewLabel(lcd, font, 0, 18, "Line2 ")
	_, labelW = label2.LineWidth()
	s.Text2 = text.NewLabel(lcd, font, int16(labelW), 18, fmt.Sprintf("%s", ""))
}

func (s *ScreenMenu) Update() {
	s.Text1.Invert = s.selectedLine == 1
	s.Text2.Invert = s.selectedLine == 2
	s.Text1.Write(fmt.Sprintf("%s", "entry1"))
	s.Text2.Write(fmt.Sprintf("%s", "entry2"))
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
