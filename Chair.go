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
	battery, speed, error uint8
	chairMsgs chan chairResponse
}

type chairResponse struct {
	typ uint8
	error uint8
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


func InitChair(c *serial.Config, stick *Joystick) Chair {
	log.Printf("Chair with path: %s", c.Name)

	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}

	chair := Chair{devicePath: c.Name, device:s, stick: stick, chairMsgs: make(chan chairResponse)}

	return chair
}

func (c *Chair) Loop() {

	stickChan := make(chan Event)
	go c.stick.readLoop(stickChan)

	go c.readLoop()

	ticker := time.Tick(25 * time.Millisecond)


	for {
		select {

		case cRes := <- c.chairMsgs:
			//log.Printf("The chair sent something: %v", cRes)
			c.battery = cRes.battery
			c.speed = cRes.speed
			c.error = cRes.error
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
			c.y = convertDirectionToChair(e.value)
		case 3: //X axis, right stick (this needs to be flipped to match the joystick)
			c.x = convertDirectionToChair(e.value) * -1
	}
}

func (c *Chair) sendData() {
	payLoad := chairData{typ: 74, command: 0, y: c.y, x: c.x}
	c.device.Write(payLoad.bytes())
}

func (d *chairData) bytes() []byte {
	bytes := []byte{d.typ, d.command, d.unknown, byte(d.x), byte(d.y), 0}
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
	fmt.Printf("\rE:%d B:%d S:%d Y:%d X:%d      ", c.error, c.battery, c.speed, c.y, c.x)
}

func (c *Chair) readLoop() {

	input := make([]byte, 5, 5)
	startByte := make([]byte, 1, 1)
	for {

		//Wait for the start byte, its 84
		for {
			c.device.Read(startByte)
			if startByte[0] == 84 {
				break;
			}
		}

		_, err := io.ReadAtLeast(c.device, input, 5)

		if err != nil {
			log.Fatal("Problem reading chair:", err)
		}

		byteReader := bytes.NewReader(input)

		cRes := chairResponse{typ: 84}

		binary.Read(byteReader, binary.LittleEndian, &cRes.error)
		binary.Read(byteReader, binary.LittleEndian, &cRes.unknown2)
		binary.Read(byteReader, binary.LittleEndian, &cRes.battery)
		binary.Read(byteReader, binary.LittleEndian, &cRes.speed)

		err = binary.Read(byteReader, binary.LittleEndian, &cRes.crc)

		if err != nil {
			log.Fatal("binary.Read failed:", err)
		}

		//log.Printf("Chair said: %v", cRes)

		c.chairMsgs <- cRes
	}
}

func convertDirectionToChair(in int16) (out int8) {
	x := ((int32(in) - 32768) * 256) / 65535
	x = x + 128

	//we are just shifting the range from -128 - 127 to -100 -100
	out = int8((x - -128) * (100 - -100) / (127 - -128) + -100)
	//log.Printf("%d became %d and then %d", in, x, out);

	return
}
