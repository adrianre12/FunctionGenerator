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
	spi *spix.SPIX
}

// New creates a new AD9833 connection. The SPI bus must already be configured.
func NewAD9833(spi *spix.SPIX) *Device {
	return &Device{
		spi: spi,
	}
}

func (d *Device) SPIwrite(tx uint16) {
	fmt.Printf("Writing: %v %x\n", tx, tx)
	_, err := d.spi.Transfer16(tx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (d *Device) SetFrequency(freq float64, register uint16) {
	register = register & (FREQ0 | FREQ1)
	freqReg := uint32(freq * math.Pow(2, 28) / 25e6)
	fmt.Printf("freqReg %x\n", freqReg)
	fmt.Printf("Low %x\n", register|uint16(freqReg&BITS14L))
	fmt.Printf("High %x\n", (register | uint16((freqReg&BITS14H)>>14)))
	d.spi.Transfer16(register | uint16(freqReg&BITS14L))
	d.spi.Transfer16(register | uint16((freqReg&BITS14H)>>14))
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
	d.SPIwrite(FREQ0 | (BITS14L & 0x10C7))
	//d.SPIwrite(0x4000)
	d.SPIwrite(FREQ0)
	//d.SPIwrite(0xC000)
	d.SPIwrite(PHASE0)
	//d.SPIwrite(0x2000)
	d.SPIwrite(B28)
}
