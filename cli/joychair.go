package main

import (
	"fmt"
	"github.com/fasmide/joychair"
)
func main() {
	fmt.Printf("Hello from main\n")
	joystick := joychair.InitJoystick("/dev/input/js1")

	fmt.Printf("X: %d, Y: %d\n", joystick.X, joystick.Y)
}
