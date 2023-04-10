package ui

import (
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenDummy struct {
	selectedLine uint16
	label1       *lcdDisplay.FieldStr
	text1        *lcdDisplay.FieldStr
	text2        *lcdDisplay.FieldStr
	text3        *lcdDisplay.FieldStr
}

func NewScreenDummy() *ScreenDummy {
	s := ScreenDummy{}
	font := &proggy.TinySZ8pt7b

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Dummy")

	return &s
}

func (s *ScreenDummy) Update() {
	lcd.WriteField(s.label1)
}

func (s *ScreenDummy) Push(result bool) {
	if result { //long push got to menu
		ChangeScreen(NewScreenMenu())
		return
	}
}

func (s *ScreenDummy) Rotate(result bool) {
	println("rotated up ", result)
}
