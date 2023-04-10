package ui

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/lcdDisplay"
	"fmt"

	"time"

	"tinygo.org/x/tinyfont/proggy"
)

var ()

type ScreenManual struct {
	label1 *lcdDisplay.FieldStr
	label2 *lcdDisplay.FieldStr
	text1  *lcdDisplay.FieldStr
	text2  *lcdDisplay.FieldStr
}

func NewScreenManual() *ScreenManual {
	s := ScreenManual{}
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Wave: ")
	_, labelW := lcd.LineWidth(s.label1)
	s.text1 = lcdDisplay.NewFieldStr(font, int16(labelW), 7, Waveform.String())

	s.label2 = lcdDisplay.NewFieldStr(font, 0, 17, "Freq: ")
	_, labelW = lcd.LineWidth(s.label2)
	s.text2 = lcdDisplay.NewFieldStr(font, int16(labelW), 17, fmt.Sprintf("%.3f", Frequency))

	return &s
}

func (s *ScreenManual) Update() {
	if Changed {
		fgen.SetMode(Waveform)
		Frequency = fgen.SetFrequency(Frequency, ad9833.ADR_FREQ0)
		Changed = false
	}

	lcd.WriteField(s.label1)
	lcd.WriteField(s.label2)

	s.text1.Value = Waveform.String()
	lcd.WriteField(s.text1)
	s.text2.Value = fmt.Sprintf("%.3f", Frequency)
	lcd.WriteField(s.text2)
}

func (s *ScreenManual) Push(result bool) {
	if result {
		ChangeScreen(NewScreenMenu())
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
