package main

import (
	ad9833 "TinyGo/FunctionGenerator/AD9833"
	text "TinyGo/FunctionGenerator/Text"
	"TinyGo/FunctionGenerator/spix"
	"TinyGo/FunctionGenerator/switches"
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/pcd8544"
	"tinygo.org/x/tinyfont/proggy"
)

var (
	spi1      = machine.SPI1
	LCDdevice *pcd8544.Device

	mode           ad9833.Mode
	frequency      float32
	changed        bool
	rotaryLastTime int64
)

var spi0 = machine.SPI0

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
	LCDdevice = pcd8544.New(spi1, dcPin, rstPin, scePin)
	LCDdevice.Configure(pcd8544.Config{})

	LCDdevice.SendCommand(pcd8544.FUNCTIONSET | pcd8544.EXTENDEDINSTRUCTION) // H = 1
	//LCDdevice.SendCommand(pcd8544.SETVOP | 0x3f)                             // 0x3f : Vop6 = 0, Vop5 to Vop0 = 1
	//LCDdevice.SendCommand(pcd8544.SETTEMP | 0x03)                            // Experimentally determined
	LCDdevice.SendCommand(pcd8544.SETBIAS | 0x04)

	LCDdevice.ClearBuffer()
	LCDdevice.ClearDisplay()
}

func SerialDelayStart(t int) {
	for i := t; i > 0; i-- {
		fmt.Printf("Starting... %d\n", i)
		time.Sleep(time.Second)
	}
}

func ConfigureRotary() {
	switches.SetupPush(machine.GP5, func(result bool) {
		if result { // ignore when pushed only use btn up
			return
		}
		fmt.Println("Released")
		mode = mode.Next()
		changed = true
		fmt.Printf("Mode %s value %v \n", mode.String(), mode.Uint16())
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

func ConfigureScreen() {
	ConfigureLCD()
	font := &proggy.TinySZ8pt7b

	label1 := text.NewLabel(LCDdevice, font, 0, 7, "Mode ")
	_, labelW := label1.LineWidth()
	text1 := text.NewLabel(LCDdevice, font, int16(labelW), 7, fmt.Sprintf("%s", mode))

	label2 := text.NewLabel(LCDdevice, font, 0, 18, "Freq ")
	_, labelW = label2.LineWidth()
	text2 := text.NewLabel(LCDdevice, font, int16(labelW), 18, fmt.Sprintf("%d", frequency))
	LCDdevice.Display()

	//Refresh screen periodically in a ro routine
	go func(t1 *text.Label, t2 *text.Label) {
		ticker := time.NewTicker(time.Millisecond * 250)

		for range ticker.C {
			if changed {
				fgen.SetMode(mode)
				frequency = fgen.SetFrequency(frequency, ad9833.ADR_FREQ0)
				changed = false
			}

			t1.Write(fmt.Sprintf("%s", mode))
			t2.Write(fmt.Sprintf("%.3f", frequency))
			LCDdevice.Display()
		}

	}(text1, text2)

}

var fgen *ad9833.Device

func main() {
	SerialDelayStart(3)
	ConfigureRotary()

	spix := spix.NewSPIX(machine.SPI0)
	spix.Configure(machine.SPIConfig{
		Frequency: 100000,
		LSBFirst:  false,
		Mode:      2,
		DataBits:  8,
		SCK:       machine.GP18,
		SDO:       machine.GP19,
	})
	spix.SetCSn(machine.GP17)
	spix.SetDatabits(16)

	fgen = ad9833.NewAD9833(spix)
	fgen.Init()
	fgen.WriteErr = false

	ConfigureScreen()

	sweepTest()

	select {}
}

func sweepTest() {
	fgen.SweepStart = 100
	fgen.SweepEnd = 10000
	fgen.SweepStep = 1
	fgen.SweepStepTime = 10
	fgen.SweepGate = true

	fgen.StartSweep()
}
