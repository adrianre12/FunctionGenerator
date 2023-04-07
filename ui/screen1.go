package ui

import (
	ad9833 "TinyGo/FunctionGenerator/AD9833"
	text "TinyGo/FunctionGenerator/Text"
	"fmt"
	"time"

	"tinygo.org/x/tinyfont/proggy"
)

var ()

type Screen1 struct {
	text1 *text.Label
	text2 *text.Label
}

func (s *Screen1) Setup() {
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	label1 := text.NewLabel(lcd, font, 0, 7, "Wave: ")
	_, labelW := label1.LineWidth()
	s.text1 = text.NewLabel(lcd, font, int16(labelW), 7, fmt.Sprintf("%s", Waveform))

	label2 := text.NewLabel(lcd, font, 0, 18, "Freq: ")
	_, labelW = label2.LineWidth()
	s.text2 = text.NewLabel(lcd, font, int16(labelW), 18, fmt.Sprintf("%f", Frequency))
}

func (s *Screen1) Update() {
	if Changed {
		fgen.SetMode(Waveform)
		Frequency = fgen.SetFrequency(Frequency, ad9833.ADR_FREQ0)
		Changed = false
	}

	s.text1.Write(fmt.Sprintf("%s", Waveform))
	s.text2.Write(fmt.Sprintf("%.3f", Frequency))
}

func (s *Screen1) Push(result bool) {
	fmt.Printf("Released %t\n", result)
	if result {
		fmt.Println("long press")
		return
	}

	Waveform = Waveform.Next()
	Changed = true
}

func (s *Screen1) Rotate(result bool) {
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
