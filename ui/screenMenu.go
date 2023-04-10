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

func NewScreenMenu() *ScreenMenu {
	s := ScreenMenu{}
	font := &proggy.TinySZ8pt7b

	s.selectedLine = 1

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Mode")
	s.text1 = lcdDisplay.NewFieldStr(font, 0, 17, "Dummy")
	s.text2 = lcdDisplay.NewFieldStr(font, 0, 27, "Manual")
	s.text3 = lcdDisplay.NewFieldStr(font, 0, 37, "Sweep")

	return &s
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
	switch s.selectedLine {
	case 1:
		{
			ChangeScreen(NewScreenDummy())
		}
	case 2:
		{
			ChangeScreen(NewScreenManual())
		}
	case 3:
		{
			ChangeScreen(NewScreenSweep())
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
