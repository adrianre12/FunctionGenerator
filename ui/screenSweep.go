package ui

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenSweep struct {
	selectedField int32
	selected      bool

	settingList *lcdDisplay.FieldList
	modeList    *lcdDisplay.FieldList
	rangeList   *lcdDisplay.FieldList
	start       *lcdDisplay.FieldFloat32
	end         *lcdDisplay.FieldFloat32
	period      *lcdDisplay.FieldInt32
	steps       *lcdDisplay.FieldInt32
}

func NewScreenSweep() *ScreenSweep {
	s := ScreenSweep{selectedField: 1}
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	s.settingList = lcdDisplay.NewFieldList(font, 0, 7, []lcdDisplay.FieldListItem{
		{Text: "Start", Value: 1},
		{Text: "End", Value: 2},
		{Text: "Period", Value: 3},
		{Text: "Steps", Value: 4},
	})

	s.rangeList = lcdDisplay.NewFieldList(font, 0, 47, []lcdDisplay.FieldListItem{
		{Text: "Hz", Value: 1},
		{Text: "KHz", Value: 1000},
	})
	s.rangeList.Selected = 1

	s.modeList = lcdDisplay.NewFieldList(font, 40, 47, []lcdDisplay.FieldListItem{
		{Text: "Sine", Value: ad9833.MODE_SINE},
		{Text: "Tri", Value: ad9833.MODE_TRI},
		{Text: "Sqr", Value: ad9833.MODE_MSB2},
	})

	s.start = lcdDisplay.NewFieldFloat32(font, 0, 37, 0)
	s.start.Format = "%.2f Hz"
	s.end = lcdDisplay.NewFieldFloat32(font, 0, 37, 0)
	s.end.Format = "%.2f Hz"
	s.period = lcdDisplay.NewFieldInt32(font, 0, 37, 10)
	s.period.Format = "%d ms/step"
	s.steps = lcdDisplay.NewFieldInt32(font, 0, 37, 10)
	s.steps.Format = "%d linear"
	return &s
}

func (s *ScreenSweep) Update() {
	if Changed {
		fgen.SetMode(uint16(s.modeList.Value()))
		//s.frequency.Value = fgen.SetFrequency(s.frequency.Value, ad9833.ADR_FREQ0)
		//Changed = false
	}

	s.settingList.Bold(s.selectedField == 1)
	s.rangeList.Bold(s.selectedField == 3)
	s.modeList.Bold(s.selectedField == 4)

	lcd.WriteField(s.settingList)
	lcd.WriteField(s.modeList)
	lcd.WriteField(s.rangeList)

	switch s.settingList.Value() {
	case 1:
		{
			s.start.Bold(s.selectedField == 2)
			lcd.WriteField(s.start)
		}
	case 2:
		{
			s.end.Bold(s.selectedField == 2)
			lcd.WriteField(s.end)
		}
	case 3:
		{
			s.period.Bold(s.selectedField == 2)
			lcd.WriteField(s.period)
		}
	case 4:
		{
			s.steps.Bold(s.selectedField == 2)
			lcd.WriteField(s.steps)
		}
	}
}

func (s *ScreenSweep) Push(result bool) {
	s.selected = !s.selected
	if result {
		ChangeScreen(NewScreenMenu())
		return
	}
}

func (s *ScreenSweep) Rotate(increment int32) {
	if s.selected {
		switch s.selectedField {
		case 1:
			{ //mode
				s.settingList.Selected = VaryBetween(s.settingList.Selected, increment, 0, 3)
				Changed = true
			}
		case 2:
			{ //setting
				switch s.settingList.Value() {
				case 1:
					{
						s.start.Value = float32(VaryBetween(int32(s.start.Value), increment*int32(s.rangeList.Value()), 0, int32(s.end.Value)-1))
					}
				case 2:
					{
						s.end.Value = float32(VaryBetween(int32(s.end.Value), increment*int32(s.rangeList.Value()), s.steps.Value+1, 2e6))
					}
				case 3:
					{
						s.period.Value = VaryBetween(s.period.Value, increment, 10, 1000)
					}
				case 4:
					{
						s.steps.Value = VaryBetween(s.steps.Value, increment, 2, 100000)
					}
				}
				s.rangeList.Selected = VaryBetween(s.rangeList.Selected, increment, 0, 1)
				Changed = true
			}
		case 3:
			{ //mode
				s.rangeList.Selected = VaryBetween(s.rangeList.Selected, increment, 0, 1)
				Changed = true
			}
		case 4:
			{ //mode
				s.modeList.Selected = VaryBetween(s.modeList.Selected, increment, 0, 2)
				Changed = true
			}
		default:
			{ // non selectable field
				s.selected = false
			}
		}
	} // this is not an if else
	if !s.selected { // not selected to scroll up an down
		s.selectedField = VaryBetween(s.selectedField, increment, 1, 4)
	}
}
