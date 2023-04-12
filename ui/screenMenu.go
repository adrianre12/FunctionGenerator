package ui

import (
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenMenu struct {
	selectedField int32
	label1        *lcdDisplay.FieldStr
	text1         *lcdDisplay.FieldStr
	text2         *lcdDisplay.FieldStr
	text3         *lcdDisplay.FieldStr
}

func NewScreenMenu() *ScreenMenu {
	s := ScreenMenu{selectedField: 1}
	font := &proggy.TinySZ8pt7b

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Mode")
	s.text1 = lcdDisplay.NewFieldStr(font, 0, 17, "Dummy")
	s.text2 = lcdDisplay.NewFieldStr(font, 0, 27, "Manual")
	s.text3 = lcdDisplay.NewFieldStr(font, 0, 37, "Sweep")

	return &s
}

func (s *ScreenMenu) Update() {
	lcd.ClearBuffer()

	lcd.WriteField(s.label1)
	s.text1.Bold(s.selectedField == 1)
	s.text2.Bold(s.selectedField == 2)
	s.text3.Bold(s.selectedField == 3)

	lcd.WriteField(s.text1)
	lcd.WriteField(s.text2)
	lcd.WriteField(s.text3)
}

func (s *ScreenMenu) Push(result bool) {
	switch s.selectedField {
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
			println("Error no screen to select for", s.selectedField)
		}

	}
}

func (s *ScreenMenu) Rotate(increment int32) {
	Changed = true
	s.selectedField = VaryBetween(s.selectedField, increment, 1, 3)
}
