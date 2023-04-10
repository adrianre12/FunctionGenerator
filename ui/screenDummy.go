package ui

import (
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenDummy struct {
	selectedField int32
	selected      bool
	label1        *lcdDisplay.FieldStr
	text1         *lcdDisplay.FieldStr
	int1          *lcdDisplay.FieldInt32
	float1        *lcdDisplay.FieldFloat32
}

func NewScreenDummy() *ScreenDummy {
	s := ScreenDummy{selectedField: 1}
	font := &proggy.TinySZ8pt7b

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Dummy")

	s.text1 = lcdDisplay.NewFieldStr(font, 0, 17, "string")
	s.int1 = lcdDisplay.NewFieldInt(font, 0, 27, 99)
	s.float1 = lcdDisplay.NewFieldFloat32(font, 0, 37, 1.0)

	return &s
}

func (s *ScreenDummy) Update() {
	s.text1.Bold(s.selectedField == 1)
	s.int1.Bold(s.selectedField == 2)
	s.float1.Bold(s.selectedField == 3)

	lcd.WriteField(s.label1)

	lcd.WriteField(s.text1)
	lcd.WriteField(s.int1)
	lcd.WriteField(s.float1)
}

func (s *ScreenDummy) Push(result bool) {
	s.selected = !s.selected
	if result { //long push got to menu
		ChangeScreen(NewScreenMenu())
		return
	}
}

func (s *ScreenDummy) Rotate(result bool) {
	println("Rotate", result)
	switch s.selectedField {
	case 2:
		{
			s.int1.Value = VaryBetween(s.int1.Value, result, 0, 100)
		}
	case 3:
		{
			s.float1.Value = float32(VaryBetween(int32(s.float1.Value*10), result, -10, 10)) / 10
		}
	default:
		{ // non selectable field
			s.selected = false
		}
	}

	if !s.selected { // not selected to scroll up an down
		s.selectedField = VaryBetween(s.selectedField, result, 1, 3)
	}

}
