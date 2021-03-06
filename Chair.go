package joychair

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"time"
)

type Chair struct {
	devicePath                            string
	device                                *serial.Port
	x, y                                  int8
	pendingCommand, battery, speed, error uint8
	chairMsgs                             chan ChairResponse
	cntr                                  uint64
	joyServer                             *JoyServer
}

type ChairResponse struct {
	typ      uint8
	error    uint8
	unknown2 uint8
	battery  uint8
	speed    uint8
	crc      uint8
}

type chairData struct {
	typ     uint8
	command uint8
	unknown uint8
	y       int8
	x       int8
	crc     uint8
}

func (d *ChairResponse) bytes() []byte {
	bytes := []byte{d.typ, d.error, d.unknown2, d.battery, d.speed, d.crc}
	return bytes
}

func InitChair(c *serial.Config) Chair {
	log.Printf("Chair with path: %s", c.Name)

	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}
	joyServer := InitJoyServer()
	chair := Chair{
		devicePath: c.Name,
		device:     s,
		chairMsgs:  make(chan ChairResponse),
		joyServer:  &joyServer}

	return chair
}

func (c *Chair) Loop() {

	go c.readLoop()

	netEventChan := make(chan JoyNetEvent)

	go c.joyServer.readLoop(netEventChan)

	ticker := time.Tick(10 * time.Millisecond)

	for {
		select {

		case cRes := <-c.chairMsgs:
			//log.Printf("The chair sent something: %v", cRes)
			c.battery = cRes.battery
			c.speed = cRes.speed
			c.error = cRes.error
			if c.cntr%5 == 1 {
				c.joyServer.send(&cRes)
			}
		case nEvent := <-netEventChan:
			c.handleNetEvent(&nEvent)
		case <-ticker:
			//log.Printf("It is time to send data to the chair")
			start := time.Now()
			c.sendData()
			c.formatCliLine(start)
		}
	}
}

func (c *Chair) handleJoystickEvent(e *Event) {
	switch e.code {
	case 2: //Y axis, right stick
		c.y = convertDirectionToChair(e.value)
	case 3: //X axis, right stick (this needs to be flipped to match the joystick)
		c.x = convertDirectionToChair(e.value) * -1
	case 8: //Speed down left trigger
		log.Printf("Slowing down %v", e.value)

		if e.value == 1 {
			c.pendingCommand = 2
		}

	case 9: //speed up, right trigger
		log.Printf("Speeding up %v", e.value)
		if e.value == 1 {
			c.pendingCommand = 4
		}
	}
}

func (c *Chair) handleNetEvent(e *JoyNetEvent) {
	c.y = e.y
	c.x = e.x
}

func (c *Chair) sendData() {
	c.cntr++
	payLoad := chairData{typ: 74, command: c.pendingCommand, y: c.y, x: c.x}

	//reset the command
	c.pendingCommand = 0

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

func (c *Chair) formatCliLine(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("\rE:%d B:%d S:%d Y:%d X:%d C:%d elpsd: %v      ", c.error, c.battery, c.speed, c.y, c.x, c.cntr, elapsed)
}

func (c *Chair) readLoop() {

	input := make([]byte, 5, 5)
	startByte := make([]byte, 1, 1)
	for {

		//Wait for the start byte, its 84
		for {
			c.device.Read(startByte)
			if startByte[0] == 84 {
				break
			}
		}

		_, err := io.ReadAtLeast(c.device, input, 5)

		if err != nil {
			log.Fatal("Problem reading chair:", err)
		}

		byteReader := bytes.NewReader(input)

		cRes := ChairResponse{typ: 84}

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

	//	return int8(in<<8)

	//we are just shifting the range from -128 - 127 to -100 -100
	out = int8((x - -128)*(100 - -100)/(127 - -128) + -100)
	//log.Printf("%d became %d and then %d", in, x, out);

	return
}
