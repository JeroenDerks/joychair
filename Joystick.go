package joychair

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"syscall"
)

type Joystick struct {
	X, Y int
	devicePath string
	device *os.File

}

type event struct {
	time syscall.Timeval
	typ, code uint16
	value uint32
}

func InitJoystick(dev string) Joystick {
	log.Printf("Joystick with path: %s", dev)

	j := Joystick{X: 0, Y:0, devicePath: dev}
	j.open()

	j.eventLoop()

	return j
}

func(j *Joystick)eventLoop() {
	for {

		input := make([]byte, 24, 24)
		j.device.Read(input)

		buf := bytes.NewReader(input)
		event := new(event)

		binary.Read(buf, binary.LittleEndian, &event.time.Sec)
		binary.Read(buf, binary.LittleEndian, &event.time.Usec)
		binary.Read(buf, binary.LittleEndian, &event.typ)
		binary.Read(buf, binary.LittleEndian, &event.code)
		err := binary.Read(buf, binary.LittleEndian, &event.value)

		if err != nil {
			log.Fatal("binary.Read failed:", err)
		}
		if event.typ == 0 && event.code == 0 && event.value == 0 {
			continue
		}
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