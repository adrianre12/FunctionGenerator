package ui

import (
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenDummy struct {
	selectedField int32
	selected      bool
	label1        *lcdDisplay.FieldStr
	list1         *lcdDisplay.FieldList
	int1          *lcdDisplay.FieldInt32
	float1        *lcdDisplay.FieldFloat32
}

func NewScreenDummy() *ScreenDummy {
	s := ScreenDummy{selectedField: 1}
	font := &proggy.TinySZ8pt7b

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Dummy")

	s.list1 = lcdDisplay.NewFieldList(font, 0, 17, []lcdDisplay.FieldListItem{
		{"One", 1},
		{"two", 2},
		{"three", 3},
	})
	s.int1 = lcdDisplay.NewFieldInt(font, 0, 27, 99)
	s.float1 = lcdDisplay.NewFieldFloat32(font, 0, 37, 1.0)

	return &s
}

func (s *ScreenDummy) Update() {
	s.list1.Bold(s.selectedField == 1)
	s.int1.Bold(s.selectedField == 2)
	s.float1.Bold(s.selectedField == 3)

	lcd.WriteField(s.label1)

	lcd.WriteField(s.list1)
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

func (s *ScreenDummy) Rotate(increment int32) {
	if s.selected {
		switch s.selectedField {
		case 1:
			{
				s.list1.Selected = VaryBetween(s.list1.Selected, increment, 0, 2)
			}
		case 2:
			{
				s.int1.Value = VaryBetween(s.int1.Value, increment, 0, 100)
			}
		case 3:
			{
				s.float1.Value = float32(VaryBetween(int32(s.float1.Value*10), increment, -10, 10)) / 10
			}
		default:
			{ // non selectable field
				s.selected = false
			}
		}
	} // this is not an if else
	if !s.selected { // not selected to scroll up an down
		s.selectedField = VaryBetween(s.selectedField, increment, 1, 3)
	}

}
