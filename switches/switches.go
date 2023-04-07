package switches

import (
	"machine"
	"time"
)

type PushButton struct {
	LongPress  uint32 // how long to hold for a long press
	Enabled    bool
	presedTime int64 // time the button was pressed in ms
}

type Callback func(result bool)

const (
	debounce = 5 //ms
)

// Creates a New PushButton tied to pin
// longPress how long the button press and release is in ms
// callback is the function called on button release, if it was a long press result = true
// The button is enabled by defult.
// func SetupPush(pin machine.Pin, longPress uint32, callback func(result bool)) (pushButton *PushButton) {
func SetupPush(pin machine.Pin, longPress uint32, callback Callback) (pushButton *PushButton) {
	pushButton = &PushButton{
		LongPress: longPress,
		Enabled:   true,
	}
	pin.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {
		if !pushButton.Enabled {
			return
		}
		now := time.Now().UnixMilli()
		if !p.Get() { // button pressed
			pushButton.presedTime = now
			return
		}

		// button released.
		if now-pushButton.presedTime < debounce { // press too short or bouncing
			return
		}

		callback(now-pushButton.presedTime >= int64(pushButton.LongPress))
	})

	return pushButton
}

// Calls back on push and release.
// result: true for pressed, active low input
func SetupButton(pin machine.Pin, callback Callback) {
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
