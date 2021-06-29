# Go NFC library example code

This is an example code for reading a NFC/RFID tag with a reader
using [NFC library for Go](https://github.com/clausecker/nfc). The NFC library requires to have libnfc 1.8.0+ installed
on your device.

I was using a Raspberry Pi with PN532 and followed the instructions from
Adafruit: [Building libnfc on Raspberry Pi](https://learn.adafruit.com/adafruit-nfc-rfid-on-raspberry-pi/building-libnfc)
.

## Building libnfc for PN532

Get and extract the libnfc:

```
 cd ~
 mkdir libnfc && cd libnfc/
 wget https://github.com/nfc-tools/libnfc/releases/download/libnfc-1.8.0/libnfc-1.8.0.tar.bz2
 tar -xvjf libnfc-1.8.0.tar.bz2
```

**Next two steps may vary for your reader**

Create PN532 configuration:

```
 cd libnfc-1.8.0
 sudo mkdir /etc/nfc
 sudo mkdir /etc/nfc/devices.d
 sudo cp contrib/libnfc/pn532_uart_on_rpi.conf.sample /etc/nfc/devices.d/pn532_uart_on_rpi.conf 
 sudo nano /etc/nfc/devices.d/pn532_uart_on_rpi.conf
```

Update the file:

> allow_intrusive_scan = true


Install dependencies for building:

```
 sudo apt-get install autoconf
 sudo apt-get install libtool
 sudo apt-get install libpcsclite-dev libusb-dev
 autoreconf -vis
 ./configure --with-drivers=pn532_uart --sysconfdir=/etc --prefix=/usr


```

Build the library:

```
sudo make clean
sudo make install all
```
