package lcdDisplay

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/pcd8544"
	"tinygo.org/x/tinyfont"
)

type Device struct {
	spi    *machine.SPI
	device *pcd8544.Device
	FG     color.RGBA
	BG     color.RGBA
}

func NewDevice() (display *Device) {
	display = &Device{
		spi: machine.SPI1,
		FG:  color.RGBA{1, 1, 1, 255},
		BG:  color.RGBA{0, 0, 0, 255},
	}

	display.spi.Configure(machine.SPIConfig{
		Frequency: 4000000,
		LSBFirst:  false,
		Mode:      0,
		DataBits:  8,
		SCK:       machine.GP10,
		SDO:       machine.GP11,
		SDI:       machine.GP12,
	})

	dcPin := machine.GP14
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin := machine.GP15
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	scePin := machine.GP9
	scePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	display.device = pcd8544.New(display.spi, dcPin, rstPin, scePin)
	display.device.Configure(pcd8544.Config{})

	display.device.SendCommand(pcd8544.FUNCTIONSET | pcd8544.EXTENDEDINSTRUCTION)
	display.device.SendCommand(pcd8544.SETBIAS | 0x04)
	display.device.ClearDisplay()
	return display
}

func (d *Device) Display() {
	d.device.Display()
}

func (d *Device) ClearBuffer() {
	d.device.ClearBuffer()
}

func (d *Device) WriteField(field Field) {
	fg := d.FG
	if field.IsBold() {
		fg = d.BG
		d.background(field, d.FG)
	}

	tinyfont.WriteLine(d.device, field.Font(), field.X(), field.Y(), field.String(), fg)
}

func (d *Device) LineWidth(field Field) (outboxWidth uint32) {
	_, outboxWidth = tinyfont.LineWidth(field.Font(), field.String())
	return outboxWidth
}

func (d *Device) background(field Field, colour color.RGBA) {
	outboxWidth := d.LineWidth(field)
	bbox := field.Font().BBox
	for x := int16(0); x < int16(outboxWidth); x++ {
		for y := int16(0); y < int16(bbox[1]); y++ {
			d.device.SetPixel(field.X()+x, field.Y()+y+int16(bbox[3]), colour)
		}
	}
}
