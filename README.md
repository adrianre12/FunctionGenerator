# FunctionGenerator
Function generator based on an AD9833 module and Pi Pico

Uses:
* Pi Pico
* AD9833 module, the most basic I could find.
* PCD8544 based 84x48 LCD. (Nokia 5110/6150 etc)
* Rotary encoder with push button

The code for manual frequency setting and for a linear sweep are functionally compleate.

TODO:
Electronics, so far I just have a few modules on a bread board.

## Connections

### AD9833
|AD9833|Pico|Pin|
|---|---|---|
|REF|||
|VCC|3V3(Out)|36|
|GND|GND|38|
|DAT|SPI0 TX|25|
|CLK|SPI0 SCK|24|
|FNC|SPI0 CSn|22|
|OUT|||

### LCD
The LCD was harvested from an old Nokia 6150, the connections below are for my wiring to the LCD. These do not match the maker modules, but the names should.
|LCD|Pico|Pin|
|---|---|---|
|1 VCC|3v3(Out)|36|
|2 SCLK|SPI1 SCK|14|
|3 MDSI|SPI1 TX|15|
|4 D/C|GP14|19|
|5 SEE|SPI1 CSn|12|
|6 GND|GND|38|
|7 RST|GP15|20|
|8 LED+|VBUS|40|

I have connected a 100 ohm resistor in series with the LEDs so they will work with the 5V VBUS and a 4.7uF capacitor beteween VOUT and GND.
For information on how to harvest your own LCD see [Here](https://community.element14.com/products/devtools/kinetiskl2freedomboard/w/documents/10856/codewarrior-tutorial-for-frdm-kl25z-zero-cost-84-48-graphical-lcd-for-the-freescale-freedom-board)

### Rotary Encoder
The rotary encoder and push button have resistor and capacitor debounce, although in code the push button code also has debounce to allow using more buttons.

Rotray connections
|Pin|Pico|Pin|
|---|---|---|
|Com|Gnd|38|
|Vcc|3V3(Out)|36|
|A|GP6|9|
|B|GP7|10|
|Sw|GP5|7|

### PWM out
During a frequency sweep a 10KHz PWM is generated repersenting the percentage of sweep. This is to be filtered and used to drive an Oscilloscope X axis.

PWM out on GP0 (pin 1)




