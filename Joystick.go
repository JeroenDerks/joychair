package joychair

import (
	"encoding/binary"
	"log"
	"os"
	"io"
	"bytes"
)

type Joystick struct {
	X, Y int
	devicePath string
	device *os.File

}

/*
 * https://www.kernel.org/doc/Documentation/input/joystick-api.txt
 * defines:
 * 	struct js_event {
 *		__u32 time;     /* event timestamp in milliseconds
 *		__s16 value;    /* value
 *		__u8 type;      /* event type
 *		__u8 number;    /* axis/button number
 *	};
 */


type event struct {
	time uint32
	value int16
	typ, code uint8
}

func InitJoystick(dev string) Joystick {
	log.Printf("Joystick with path: %s", dev)

	j := Joystick{X: 0, Y:0, devicePath: dev}
	j.open()

	j.eventLoop()

	return j
}

func (j *Joystick) eventLoop() {

	input := make([]byte, 24, 24)

	for {



		io.ReadAtLeast(j.device, input, 24)

		buf := bytes.NewReader(input)

		event := new(event)

		binary.Read(buf, binary.LittleEndian, &event.time)
		binary.Read(buf, binary.LittleEndian, &event.value)
		binary.Read(buf, binary.LittleEndian, &event.typ)
		err := binary.Read(buf, binary.LittleEndian, &event.code)

		if err != nil {
			log.Fatal("binary.Read failed:", err)
		}
		// if event.typ == 0 && event.code == 0 && event.value == 0 {
		// 	continue
		// }
		//... so what now?
		log.Printf("I had a event: %v", event);

	}
}

func (j *Joystick)open() error {
	file, err := os.Open(j.devicePath)

	if err != nil {
		log.Fatal(err)
	}

	j.device = file

	return nil
}