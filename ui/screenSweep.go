package ui

import (
	"TinyGo/FunctionGenerator/lcdDisplay"
)

type ScreenSweep struct {
	selectedLine uint16
	text1        *lcdDisplay.FieldStr
	text2        *lcdDisplay.FieldStr
}

func NewScreenSweep() *ScreenSweep {
	s := ScreenSweep{}

	return &s
}

func (s *ScreenSweep) Update() {
	panic("not implemented") // TODO: Implement
}

func (s *ScreenSweep) Push(result bool) {
	panic("not implemented") // TODO: Implement
}

func (s *ScreenSweep) Rotate(result bool) {
	panic("not implemented") // TODO: Implement
}
