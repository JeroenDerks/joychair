package joychair

import (
	"encoding/binary"
	"log"
	"os"
	"bytes"
)

type Joystick struct {
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


type Event struct {
	time uint32
	value int16
	typ, code uint8
}

func InitJoystick(dev string) Joystick {
	log.Printf("Joystick with path: %s", dev)

	j := Joystick{devicePath: dev}
	j.open()

	return j
}

func (j *Joystick) readLoop(c chan Event) {

	input := make([]byte, 8, 8)

	for {

		_, err := j.device.Read(input)

		if err != nil {
			log.Fatal("Problem reading joystick:", err)
		}

		byteReader := bytes.NewReader(input)

		event := new(Event)

		binary.Read(byteReader, binary.LittleEndian, &event.time)
		binary.Read(byteReader, binary.LittleEndian, &event.value)
		binary.Read(byteReader, binary.LittleEndian, &event.typ)

		err = binary.Read(byteReader, binary.LittleEndian, &event.code)

		if err != nil {
			log.Fatal("binary.Read failed:", err)
		}

		c <- *event
	}
}

func (j *Joystick) open() error {
	file, err := os.Open(j.devicePath)

	if err != nil {
		log.Fatal(err)
	}

	j.device = file

	return nil
}