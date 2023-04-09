package ui

import (
	"TinyGo/FunctionGenerator/lcdDisplay"
)

type ScreenSweep struct {
	selectedLine uint16
	text1        lcdDisplay.Field
	text2        lcdDisplay.Field
}

func (s *ScreenSweep) Setup() {
	panic("not implemented") // TODO: Implement
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
