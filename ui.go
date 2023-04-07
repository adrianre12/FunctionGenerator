package main

import (
	"TinyGo/FunctionGenerator/switches"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/pcd8544"
)

type Display interface {
	drivers.Displayer
	ClearBuffer()
	ClearDisplay()
}

type Screen interface {
	Setup()
	Update()
	Push(result bool)
	Rotate(result bool)
}

var (
	lcd Display

	changed        bool
	rotaryLastTime int64
	screen         Screen
	screen1        Screen
)

func ConfigureLCD() {
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
	lcdDevice := pcd8544.New(spi1, dcPin, rstPin, scePin)
	lcdDevice.Configure(pcd8544.Config{})

	lcdDevice.SendCommand(pcd8544.FUNCTIONSET | pcd8544.EXTENDEDINSTRUCTION)
	lcdDevice.SendCommand(pcd8544.SETBIAS | 0x04)
	lcd = lcdDevice //use Device interface for the rest so the lcd device can be changed.
	lcd.ClearDisplay()
}

func ConfigureScreen() {
	screen1 = &Screen1{}
	screen = screen1
	ConfigureLCD()

	screen.Setup()
	lcd.Display()

	//Refresh screen periodically in a go routine
	//go func(t1 *text.Label, t2 *text.Label) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 250)

		for range ticker.C {
			screen.Update()
			lcd.Display()
		}

	}()

}

func ConfigureKeyboard() {
	switches.SetupPush(machine.GP5, 1000, func(result bool) {
		screen.Push(result) //have to do it via this func to avoid runtime panic
	})

	switches.SetupRotary(machine.GP6, machine.GP7, func(result bool) {
		screen.Rotate(result) //have to do it via this func to avoid runtime panic
	})
}
