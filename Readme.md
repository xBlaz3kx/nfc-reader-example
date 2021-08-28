# Go NFC library example code

This is an example code for reading a NFC/RFID tag with a reader
using [NFC library for Go](https://github.com/clausecker/nfc). The NFC library requires to have libnfc 1.8.0+ installed
on your device.

I was using a Raspberry Pi with PN532 and followed the instructions from
Adafruit: [Building libnfc on Raspberry Pi](https://learn.adafruit.com/adafruit-nfc-rfid-on-raspberry-pi/building-libnfc)
.

## Building libnfc for PN532

1. Get and extract the libnfc:

    ```bash
     cd ~
     mkdir libnfc && cd libnfc/
     wget https://github.com/nfc-tools/libnfc/releases/download/libnfc-1.8.0/libnfc-1.8.0.tar.bz2
     tar -xvjf libnfc-1.8.0.tar.bz2
    ```

   **Next two steps depend on your reader and it's configuration of libnfc**

2. Create PN532 configuration:

    ```bash
     cd libnfc-1.8.0
     sudo mkdir /etc/nfc
     sudo mkdir /etc/nfc/devices.d
     sudo cp contrib/libnfc/pn532_uart_on_rpi.conf.sample /etc/nfc/devices.d/pn532_uart_on_rpi.conf 
     sudo nano /etc/nfc/devices.d/pn532_uart_on_rpi.conf
    ```

3. Update the _pn532_uart_on_rpi.conf_:

   ```text
   allow_intrusive_scan = true
   ```

4. Install dependencies for building:

   ```bash
    sudo apt-get install autoconf
    sudo apt-get install libtool
    sudo apt-get install libpcsclite-dev libusb-dev
    autoreconf -vis
    ./configure --with-drivers=pn532_uart --sysconfdir=/etc --prefix=/usr
   ```

5. Build the library:

   ```bash
   sudo make clean
   sudo make install all
   ```

## Running the code:

Clone the repo:

   ```bash
   git clone https://github.com/xBlaz3kx/nfc-reader-go-example/
   ```

Install the libnfc, then run:

   ```bash
   go run nfc-reader-example.go 
   ```

or

  ```bash
   go build nfc-reader-example.go && ./nfc-reader-example
   ```
