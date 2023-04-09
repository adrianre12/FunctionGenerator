package ui

import (
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenMenu struct {
	selectedLine uint16
	label1       lcdDisplay.Field
	text1        lcdDisplay.Field
	text2        lcdDisplay.Field
	text3        lcdDisplay.Field
}

func (s *ScreenMenu) Setup() {
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	s.selectedLine = 1

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Mode")
	s.text1 = lcdDisplay.NewFieldStr(font, 0, 17, "Dummy")
	s.text2 = lcdDisplay.NewFieldStr(font, 0, 27, "Manual")
	s.text3 = lcdDisplay.NewFieldStr(font, 0, 37, "Sweep")
}

func (s *ScreenMenu) Update() {
	//s.Label1.Bold(true)
	lcd.WriteField(s.label1)
	s.text1.Bold(s.selectedLine == 1)
	s.text2.Bold(s.selectedLine == 2)
	s.text3.Bold(s.selectedLine == 3)

	lcd.WriteField(s.text1)
	lcd.WriteField(s.text2)
	lcd.WriteField(s.text3)
}

func (s *ScreenMenu) Push(result bool) {
	println("Selected line", s.selectedLine)
	switch s.selectedLine {
	case 1:
		{
			println("Select Dummy")
			ChangeScreen(Dummy)
		}
	case 2:
		{
			println("select Manual")
			ChangeScreen(Manual)
		}
	case 3:
		{
			ChangeScreen(Sweep)
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
	if !result && s.selectedLine < 3 { //down
		s.selectedLine++
	}
}
