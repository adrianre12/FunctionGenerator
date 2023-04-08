package ui

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/text"
	"fmt"

	"time"

	"tinygo.org/x/tinyfont/proggy"
)

var ()

type ScreenManual struct {
	text1 *text.Label
	text2 *text.Label
}

func (s *ScreenManual) Setup() {
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	label1 := text.NewLabel(lcd, font, 0, 7, "Wave: ")
	_, labelW := label1.LineWidth()
	s.text1 = text.NewLabel(lcd, font, int16(labelW), 7, Waveform.String())

	label2 := text.NewLabel(lcd, font, 0, 17, "Freq: ")
	_, labelW = label2.LineWidth()
	s.text2 = text.NewLabel(lcd, font, int16(labelW), 17, fmt.Sprintf("%.3f", Frequency))
}

func (s *ScreenManual) Update() {
	if Changed {
		fgen.SetMode(Waveform)
		Frequency = fgen.SetFrequency(Frequency, ad9833.ADR_FREQ0)
		Changed = false
	}

	s.text1.Write(Waveform.String())
	s.text2.Write(fmt.Sprintf("%.3f", Frequency))
}

func (s *ScreenManual) Push(result bool) {
	if result {
		ChangeScreen(Menu)
		return
	}

	Waveform = Waveform.Next()
	Changed = true
}

func (s *ScreenManual) Rotate(result bool) {
	delta := time.Now().UnixMilli() - rotaryLastTime
	var increment float32
	switch {
	case delta < 25:
		{
			increment = 50
		}
	case delta < 75:
		{
			increment = 10
		}
	case delta < 150:
		{
			increment = 5
		}
	default:
		{
			increment = 1
		}
	}
	if result {
		Frequency += increment
	} else {
		Frequency -= increment
	}
	if Frequency < 0 {
		Frequency = 0
	}
	Changed = true
	rotaryLastTime = time.Now().UnixMilli()
}
