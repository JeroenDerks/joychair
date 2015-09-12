package joychair

import (
	"log"
	"bytes"
	"encoding/binary"
	"github.com/tarm/serial"
	"time"
	"io"
	"fmt"
)

type Chair struct {
	devicePath string
	device *serial.Port
	stick *Joystick
	x, y int8
	battery, speed uint8
	chairMsgs chan chairResponse
}

type chairResponse struct {
	typ uint8
	unknown uint8
	unknown2 uint8
	battery uint8
	speed uint8
	crc uint8
}

type chairData struct {
	typ uint8
	command uint8
	unknown uint8
	y int8
	x int8
	crc uint8
}


func InitChair(dev string, stick *Joystick) Chair {
	log.Printf("Chair with path: %s", dev)

	c := &serial.Config{Name: dev, Baud: 115200}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}

	chair := Chair{devicePath: dev, device:s, stick: stick, chairMsgs: make(chan chairResponse)}

	return chair
}

func (c *Chair) Loop() {

	stickChan := make(chan Event)
	go c.stick.readLoop(stickChan)

	go c.readLoop()

	ticker := time.Tick(15 * time.Millisecond)


	for {
		select {

		case cRes := <- c.chairMsgs:
			//log.Printf("The chair sent something: %v", cRes)
			c.battery = cRes.battery
			c.speed = cRes.speed
			c.formatCliLine()
		case sEvent := <- stickChan:
			//log.Printf("Stick sent something: %v", sEvent)
			c.handleJoystickEvent(&sEvent)
			c.formatCliLine()
		case <- ticker:
			//log.Printf("It is time to send data to the chair")
			c.sendData()
		}
	}
}

func (c *Chair) handleJoystickEvent(e *Event) {
	switch e.code {
		case 2: //Y axis, right stick
			c.y = int8(e.value)
		case 3: //X axis, right stick
			c.x = int8(e.value)
	}
}

func (c *Chair) sendData() {
	payLoad := chairData{typ: 74, command: 0, y: c.y, x: c.x}
	c.device.Write(payLoad.bytes())
}

func (d *chairData) bytes() []byte {
	bytes := []byte{d.typ, d.command, d.unknown, byte(d.y), byte(d.x), 0}
	bytes[5] = calculateCheckSum(bytes)
	return bytes
}

func calculateCheckSum(b []byte) byte {
	sum := byte(255)

	for i := 0; i < 5; i++ {
		sum = sum - b[i]
	}

	return sum
}

func (c *Chair) formatCliLine() {
	fmt.Printf("\rB:%d S:%d Y:%d X:%d      ", c.battery, c.speed, c.y, c.x)
}

func (c *Chair) readLoop() {

	input := make([]byte, 6, 6)

	for {

		_, err := io.ReadAtLeast(c.device, input, 6)

		if err != nil {
			log.Fatal("Problem reading chair:", err)
		}

		byteReader := bytes.NewReader(input)

		cRes := new(chairResponse)

		binary.Read(byteReader, binary.LittleEndian, &cRes.typ)
		binary.Read(byteReader, binary.LittleEndian, &cRes.unknown)
		binary.Read(byteReader, binary.LittleEndian, &cRes.unknown2)
		binary.Read(byteReader, binary.LittleEndian, &cRes.battery)
		binary.Read(byteReader, binary.LittleEndian, &cRes.speed)

		err = binary.Read(byteReader, binary.LittleEndian, &cRes.crc)

		if err != nil {
			log.Fatal("binary.Read failed:", err)
		}

		c.chairMsgs <- *cRes
	}
}