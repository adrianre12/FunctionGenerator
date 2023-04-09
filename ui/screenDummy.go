package ui

import (
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenDummy struct {
	selectedLine uint16
	label1       lcdDisplay.Field
	text1        lcdDisplay.Field
	text2        lcdDisplay.Field
	text3        lcdDisplay.Field
}

func (s *ScreenDummy) Setup() {
	println("ScreenDummy")
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Dummy")
}

func (s *ScreenDummy) Update() {
	panic("not implemented") // TODO: Implement
}

func (s *ScreenDummy) Push(result bool) {
	panic("not implemented") // TODO: Implement
}

func (s *ScreenDummy) Rotate(result bool) {
	panic("not implemented") // TODO: Implement
}
