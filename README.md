# Joychair
A remote controlled wheelchair ... ish

It uses an raspberry pi, ps3 controller and a old wheelchair with an controller marked VR2. (http://sunrise.pgdrivestechnology.com/pdf/vr2.pdf)

# Setting up a raspberry pi

## Install dependencies
```
sudo apt-get update

sudo apt-get install bluez-utils bluez-compat bluez-hcidump checkinstall libusb-dev libbluetooth-dev joystick

wget http://www.pabr.org/sixlinux/sixpair.c
gcc -o sixpair sixpair.c -lusb
```

Download qtsixa from http://qtsixa.sourceforge.net/
```
tar xfvz QtSixA-1.5.1-src.tar.gz
```
At this point you may need to reboot the pi. I was unable to compile the sixad binary due to "virtual memory exhausted: Cannot allocate memory" errors

```
cd QtSixA-1.5.1/sixad
make
sudo make install
```

## Pairing

Connect the controller via the USB bus
```
sudo ~/sixpair
```
Disconnect usb and run:
```
sudo sixad -s
```
Test with:
```
jstest /dev/input/js0
```

## Starting up on boot
Ive had some trouble getting the ps3 controller and pi reliably connect to each other. I am doing this to make it kind of reliable:

1. make sure the bluetooth service is stopped (`service bluetooth status` should tell you its stopped)
2. When the pi boots up, run `hciconfig hci0 up pscan`
3. sixad-bin is evil, when it stops it starts the bluetooth service again - making the pi and controller unable to communicate - so make sure you stop it.

Start off by making sure bluetooth is not started on boot
```
sudo update-rc.d -f bluetooth remove
```

Then, put this at the *end* of your `/etc/rc.local`

```
hciconfig hci0 up
hciconfig hci0 pscan
sixad-bin 0 0 0 &
```

As you might have guessed, if sixad-bin decides to exit or crash, it may startup the bluetooth service once again and your out of luck :)
