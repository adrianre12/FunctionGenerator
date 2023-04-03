package switches

import "machine"

type Callback func(result bool)

// Calls back on push and release.
//result: true for pressed, active low input
func SetupPush(pin machine.Pin, callback Callback) {
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	pin.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {
		callback(!p.Get())
	})
}

// Calls back on rotation.
// result: true and false indicate direction
func SetupRotary(pinA machine.Pin, pinB machine.Pin, callback Callback) {
	pinA.Configure(machine.PinConfig{Mode: machine.PinInput})
	pinB.Configure(machine.PinConfig{Mode: machine.PinInput})
	pinA.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		callback(p.Get() == pinB.Get())
	})
}
