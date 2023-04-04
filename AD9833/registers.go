package ad9833

import "fmt"

const (
	PHASE1  uint16 = 0xE000 //D15/D14/D13 0xE000
	PHASE0  uint16 = 0xC000 //D15/D14/D13 0xC000
	FREQ1   uint16 = 0x8000 //D15/D14 0x8000
	FREQ0   uint16 = 0x4000 //D15/D14 0x4000
	B28     uint16 = 0x2000 //D13 0x2000
	HLB     uint16 = 0x1000 //D12 0x1000
	FSELECT uint16 = 0x0800 //D11 0x0800
	PSELECT uint16 = 0x0400 //D10 0x0400
	//		= 0x0200 //D9 Reserved
	RESET   uint16 = 0x0100 //D8 0x0100
	SLEEP1  uint16 = 0x0080 //D7 0x0080
	SLEEP12 uint16 = 0x0040 //D6 0x0040
	OPBITEN uint16 = 0x0020 //D5 0x0020
	//		= 0x0010 //D4 Reserved
	DIV2 uint16 = 0x0008 //D3 0x0008
	//		= 0x0004 //D2 Reserved
	MODE uint16 = 0x0002 //D1 0x0002
	//		= 0x0001 //D0 Reserved

	BITS14H uint32 = 0x0FFFC000 //High 14 MSBs 0xFFFC000
	BITS14L uint32 = 0x00003FFF //Low 14 LSBs 0x3FFF
	BITS12  uint16 = 0x0FFF     //12 LSBs 0x0FFF
)

type Mode uint8

const (
	MODE_MASK uint16 = 0x002A //D5/D3/D1 Mask 0x002A
)
const (
	MODE_SINE Mode = iota //D5/D3/D1 Sinewave 0x0000
	MODE_TRI              //D5/D3/D1 Triangle 0x0008
	MODE_MSB2             //D5/D3/D1 MSB/2 0x0020
	MODE_MSB              //D5/D3/D1 MSB 0x0020
)

func (m Mode) String() (str string) {
	switch m {
	case MODE_SINE:
		{
			str = "Sine"
		}
	case MODE_TRI:
		{
			str = "Tri"
		}
	case MODE_MSB2:
		{
			str = "Sqr"
		}
	case MODE_MSB:
		{
			str = "Sqr*2"
		}
	}
	return str
}

func (m Mode) Uint16() (value uint16) {
	switch m {
	case MODE_SINE:
		{
			value = 0x0000 //D5/D3/D1 Sinewave 0x0000
		}
	case MODE_TRI:
		{
			value = 0x0002 //D5/D3/D1 Triangle 0x0008
		}
	case MODE_MSB2:
		{
			value = 0x0020 //D5/D3/D1 MSB/2 0x0020
		}
	case MODE_MSB:
		{
			value = 0x0022 //D5/D3/D1 MSB 0x0020
		}
	}
	return value
}

func (m Mode) Next() (next Mode) {
	next = m + 1
	fmt.Printf("m=%d next=%d", uint8(m), uint8(next))

	if m >= 2 {
		next = 0
	}
	fmt.Printf(" Set=%d\n", uint8(next))
	return next
}

type register struct {
	value uint16
}

func (r *register) replaceBits(value uint16, mask uint16) {
	r.value = (r.value & ^mask) | (value & mask)
}
