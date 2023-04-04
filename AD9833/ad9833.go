package ad9833

import (
	"TinyGo/FunctionGenerator/spix"
	"errors"
	"fmt"
	"math"
)

const SPITimeout = 10000 //10ms
var (
	ErrSPITimeout = errors.New("SPI timeout")
)

type Device struct {
	spi        *spix.SPIX
	controlReg register
}

// New creates a new AD9833 connection. The SPI bus must already be configured.
func NewAD9833(spi *spix.SPIX) *Device {
	return &Device{
		spi:        spi,
		controlReg: register{value: 0},
	}
}

func (d *Device) Init() {
	d.controlReg.value = uint16(B28 | RESET)
	d.spi.Transfer16(d.controlReg.value)

	//set freq and phase to 0
	d.spi.Transfer16(FREQ0) //LSB
	d.spi.Transfer16(FREQ0) //MSB

	d.spi.Transfer16(FREQ1) //LSB
	d.spi.Transfer16(FREQ1) //MSB

	d.spi.Transfer16(PHASE0)
	d.spi.Transfer16(PHASE1)

	d.controlReg.replaceBits(0, RESET)
	d.spi.Transfer16(d.controlReg.value)
}

func (d *Device) SPIwrite(tx uint16) {
	//fmt.Printf("Writing: %v %x\n", tx, tx)
	_, err := d.spi.Transfer16(tx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (d *Device) SetFrequency(freq float64, freqReg uint16) {
	d.spi.Transfer16(d.controlReg.value)
	freqReg = freqReg & (FREQ0 | FREQ1)
	freqValue := uint32(freq * math.Pow(2, 28) / 25e6)
	//fmt.Printf("freqReg %x\n", freqReg)
	//fmt.Printf("Low %x\n", freqReg|uint16(freqValue&BITS14L))
	//fmt.Printf("High %x\n", (freqReg | uint16((freqValue&BITS14H)>>14)))
	d.spi.Transfer16(freqReg | uint16(freqValue&BITS14L))
	d.spi.Transfer16(freqReg | uint16((freqValue&BITS14H)>>14))
}

func (d *Device) SetMode(mode Mode) {
	d.controlReg.replaceBits(mode.Uint16(), MODE_MASK)
	d.spi.Transfer16(d.controlReg.value)
}

func (d *Device) FreqTest() {
	/* From document AN-1070
	   0x2100 0010 0001 0000 0000
	   0x50C7 0101 0000 1100 0111
	   0x4000 0100 0000 0000 0000
	   0xC000 1100 0000 0000 0000
	   0x2000 0010 0000 0000 0000
	*/
	fmt.Println("freqTest")
	//d.SPIwrite(0x2100)
	d.SPIwrite(B28 | RESET)
	//d.SPIwrite(0x50C7)
	d.SPIwrite(FREQ0 | uint16(BITS14L&0x10C7))
	//d.SPIwrite(0x4000)
	d.SPIwrite(FREQ0)
	//d.SPIwrite(0xC000)
	d.SPIwrite(PHASE0)
	//d.SPIwrite(0x2000)
	d.SPIwrite(B28)
}
