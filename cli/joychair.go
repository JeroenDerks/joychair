package main

import (
	"fmt"
	"github.com/fasmide/joychair"
)


func main() {
	fmt.Printf("Hello from main\n")
	joystick := joychair.InitJoystick("/dev/input/js1")

	chair := joychair.InitChair("/dev/ttyACM0", &joystick)

	chair.Loop()


	fmt.Printf("No more waiting")

}
