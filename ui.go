package main

import (
	ad9833 "TinyGo/FunctionGenerator/AD9833"
	text "TinyGo/FunctionGenerator/Text"
	"TinyGo/FunctionGenerator/switches"
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/pcd8544"
	"tinygo.org/x/tinyfont/proggy"
)

var (
	lcd *pcd8544.Device

	changed        bool
	rotaryLastTime int64
)

func configureLCD() {
	spi1.Configure(machine.SPIConfig{
		Frequency: 4000000,
		LSBFirst:  false,
		Mode:      0,
		DataBits:  8,
		SCK:       machine.GP10,
		SDO:       machine.GP11,
	})

	dcPin := machine.GP14
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin := machine.GP15
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	scePin := machine.GP9
	scePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	lcd = pcd8544.New(spi1, dcPin, rstPin, scePin)
	lcd.Configure(pcd8544.Config{})

	lcd.SendCommand(pcd8544.FUNCTIONSET | pcd8544.EXTENDEDINSTRUCTION) // H = 1
	//LCDdevice.SendCommand(pcd8544.SETVOP | 0x3f)                             // 0x3f : Vop6 = 0, Vop5 to Vop0 = 1
	//LCDdevice.SendCommand(pcd8544.SETTEMP | 0x03)                            // Experimentally determined
	lcd.SendCommand(pcd8544.SETBIAS | 0x04)

	lcd.ClearBuffer()
	lcd.ClearDisplay()
}

func ConfigureScreen() {
	configureLCD()
	font := &proggy.TinySZ8pt7b

	label1 := text.NewLabel(lcd, font, 0, 7, "Wave ")
	_, labelW := label1.LineWidth()
	text1 := text.NewLabel(lcd, font, int16(labelW), 7, fmt.Sprintf("%s", waveform))

	label2 := text.NewLabel(lcd, font, 0, 18, "Freq ")
	_, labelW = label2.LineWidth()
	text2 := text.NewLabel(lcd, font, int16(labelW), 18, fmt.Sprintf("%f", frequency))
	lcd.Display()

	//Refresh screen periodically in a go routine
	go func(t1 *text.Label, t2 *text.Label) {
		ticker := time.NewTicker(time.Millisecond * 250)

		for range ticker.C {
			if changed {
				fgen.SetMode(waveform)
				frequency = fgen.SetFrequency(frequency, ad9833.ADR_FREQ0)
				changed = false
			}

			t1.Write(fmt.Sprintf("%s", waveform))
			t2.Write(fmt.Sprintf("%.3f", frequency))
			lcd.Display()
		}

	}(text1, text2)

}

func ConfigureRotary() {
	switches.SetupPush(machine.GP5, 1000, func(result bool) {
		fmt.Printf("Released %t\n", result)
		if result {
			fmt.Println("long press")
			return
		}

		waveform = waveform.Next()
		changed = true
	})

	switches.SetupRotary(machine.GP6, machine.GP7, func(result bool) {
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
			frequency += increment
		} else {
			frequency -= increment
		}
		if frequency < 0 {
			frequency = 0
		}
		changed = true
		rotaryLastTime = time.Now().UnixMilli()
	})
}
