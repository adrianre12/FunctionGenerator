package main

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/spix"
	"TinyGo/FunctionGenerator/sweep"
	"TinyGo/FunctionGenerator/ui"

	"fmt"
	"machine"
	"time"
)

func SerialDelayStart(t int) {
	for i := t; i > 0; i-- {
		fmt.Printf("Starting... %d\n", i)
		time.Sleep(time.Second)
	}
}

func main() {
	SerialDelayStart(5)

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

	fgen := ad9833.NewAD9833(spix)
	fgen.Init()
	fgen.WriteErr = false

	sweep.ConfigureSweep(fgen)
	sweep.SetPWM(0.5)
	ui.Configure(fgen) //ui should be last

	select {}
}
