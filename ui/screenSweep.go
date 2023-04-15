package ui

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/lcdDisplay"
	"TinyGo/FunctionGenerator/sweep"
	"time"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenSweep struct {
	selectedField int32
	selected      bool

	settingList *lcdDisplay.FieldList
	modeList    *lcdDisplay.FieldList
	rangeList   *lcdDisplay.FieldList
	start       *lcdDisplay.FieldFloat64
	end         *lcdDisplay.FieldFloat64
	stepTime    *lcdDisplay.FieldInt32
	steps       *lcdDisplay.FieldInt32
	run         *lcdDisplay.FieldStr
	sweeping    *lcdDisplay.FieldStr
}

var sweeping bool

func NewScreenSweep() *ScreenSweep {
	s := ScreenSweep{selectedField: 1}
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	s.modeList = lcdDisplay.NewFieldList(font, 0, 7, []lcdDisplay.FieldListItem{
		{Text: "Sine", Value: ad9833.MODE_SINE},
		{Text: "Tri", Value: ad9833.MODE_TRI},
		{Text: "Sqr", Value: ad9833.MODE_MSB2},
	})

	s.settingList = lcdDisplay.NewFieldList(font, 0, 17, []lcdDisplay.FieldListItem{
		{Text: "Start", Value: 1},
		{Text: "End", Value: 2},
		{Text: "Step Time", Value: 3},
		{Text: "Steps", Value: 4},
	})

	s.start = lcdDisplay.NewFieldFloat64(font, 0, 37, 10)
	s.start.Format = "%.0f Hz"
	s.end = lcdDisplay.NewFieldFloat64(font, 0, 37, 5000)
	s.end.Format = "%.0f Hz"
	s.stepTime = lcdDisplay.NewFieldInt32(font, 0, 37, 20)
	s.stepTime.Format = "%d ms/step"
	s.steps = lcdDisplay.NewFieldInt32(font, 0, 37, 1000)
	s.steps.Format = "%d linear"

	s.rangeList = lcdDisplay.NewFieldList(font, 0, 47, []lcdDisplay.FieldListItem{
		{Text: "Hz", Value: 1},
		{Text: "KHz", Value: 1000},
	})
	s.rangeList.Selected = 1

	s.run = lcdDisplay.NewFieldStr(font, 40, 47, "Run")

	s.sweeping = lcdDisplay.NewFieldStr(font, 15, 20, "Sweeping")

	return &s
}

func (s *ScreenSweep) Update() {
	lcd.ClearBuffer()

	if sweeping {
		lcd.WriteField(s.sweeping)
		return
	}

	s.modeList.Bold(s.selectedField == 1)
	s.settingList.Bold(s.selectedField == 2)
	s.rangeList.Bold(s.selectedField == 4)

	switch s.settingList.Value() {
	case 1:
		{
			s.start.Bold(s.selectedField == 3)
			lcd.WriteField(s.start)
		}
	case 2:
		{
			s.end.Bold(s.selectedField == 3)
			lcd.WriteField(s.end)
		}
	case 3:
		{
			s.stepTime.Bold(s.selectedField == 3)
			lcd.WriteField(s.stepTime)
		}
	case 4:
		{
			s.steps.Bold(s.selectedField == 3)
			lcd.WriteField(s.steps)
		}
	}

	s.run.Bold(s.selectedField == 5)

	lcd.WriteField(s.settingList)
	lcd.WriteField(s.modeList)
	lcd.WriteField(s.rangeList)
	lcd.WriteField(s.run)

}

func (s *ScreenSweep) Push(result bool) {
	s.selected = !s.selected
	if result {
		ChangeScreen(NewScreenMenu())
		return
	}

	switch s.selectedField {
	case 5:
		{ //run
			Changed = true
			if sweeping {
				sweep.StopSweep()
			} else {
				go sweep.StartSweep(uint16(s.modeList.Value()), s.start.Value, s.end.Value, s.stepTime.Value, s.steps.Value, true, func() {
					sweeping = false
					s.selected = false
				})
			}
			sweeping = s.selected
			if sweeping {
				go func() {
					var bold bool
					for sweeping {
						bold = !bold
						s.sweeping.Bold(bold)
						time.Sleep(time.Second)
						Changed = true
					}
				}()
			}
		}
	}
}

func (s *ScreenSweep) Rotate(increment int32) {
	Changed = true
	if s.selected {
		switch s.selectedField {
		case 1:
			{ //mode
				s.modeList.Selected = VaryInt32Between(s.modeList.Selected, increment, 0, 2)
			}
		case 2:
			{ //setting
				s.settingList.Selected = VaryInt32Between(s.settingList.Selected, increment, 0, 3)
			}
		case 3:
			{ //setting value fields
				switch s.settingList.Value() {
				case 1:
					{
						s.start.Value = VaryFloat64Between(s.start.Value, float64(increment)*float64(s.rangeList.Value()), 0, float64(s.end.Value)-1)
					}
				case 2:
					{
						s.end.Value = VaryFloat64Between(s.end.Value, float64(increment)*float64(s.rangeList.Value()), float64(s.steps.Value+1), 2e6)
					}
				case 3:
					{
						s.stepTime.Value = VaryInt32Between(s.stepTime.Value, increment, 10, 1000)
					}
				case 4:
					{
						s.steps.Value = VaryInt32Between(s.steps.Value, increment, 2, 100000)
					}
				}
			}
		case 4:
			{ //range
				s.rangeList.Selected = VaryInt32Between(s.rangeList.Selected, increment, 0, 1)
			}
		case 5:
			{ // run
				//discard rotation while run is selected.
			}
		default:
			{ // non selectable field
				s.selected = false
			}
		}
	} // this is not an if else
	if !s.selected { // not selected to scroll up an down
		s.selectedField = VaryInt32Between(s.selectedField, increment, 1, 5)
	}
}
