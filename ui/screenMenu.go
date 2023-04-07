package ui

import (
	text "TinyGo/FunctionGenerator/Text"
	"fmt"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenMenu struct {
	Text1 *text.Label
	Text2 *text.Label
}

func (s *ScreenMenu) Setup() {
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	label1 := text.NewLabel(lcd, font, 0, 7, "Line1 ")
	_, labelW := label1.LineWidth()
	s.Text1 = text.NewLabel(lcd, font, int16(labelW), 7, fmt.Sprintf("%s", ""))

	label2 := text.NewLabel(lcd, font, 0, 18, "Line2 ")
	_, labelW = label2.LineWidth()
	s.Text2 = text.NewLabel(lcd, font, int16(labelW), 18, fmt.Sprintf("%s", ""))
	s.Text2.Invert = true
}

func (s *ScreenMenu) Update() {
	s.Text1.Write(fmt.Sprintf("%s", "entry1"))
	s.Text2.Write(fmt.Sprintf("%s", "entry2"))
}

func (s *ScreenMenu) Push(result bool) {
	println("Push", result)
}

func (s *ScreenMenu) Rotate(result bool) {
	println("Rotate", result)
}
