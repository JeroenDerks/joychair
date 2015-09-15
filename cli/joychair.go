package main

import (
	"os"
	"log"
	"github.com/tarm/serial"
	"github.com/fasmide/joychair"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Chair serial.Config
	Joystick joychair.JoystickConfig
}

func main() {
	
	if len(os.Args) <= 1 {
		log.Fatal("Provide configuration path as first argument");
	}

	configPath := os.Args[1]

	config := Config{}

	log.Printf("Reading configuration from %s", configPath)

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal(err)
	}

	joystick := joychair.InitJoystick(&config.Joystick)

	chair := joychair.InitChair(&config.Chair, &joystick)

	chair.Loop()


	log.Printf("Bye")

}
