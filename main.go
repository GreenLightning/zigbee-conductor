package main

import (
	"fmt"

	"github.com/jacobsa/go-serial/serial"
)

func main() {
	options := serial.OpenOptions{
		PortName:          "/dev/ttyACM0",
		BaudRate:          115200,
		RTSCTSFlowControl: true,
		DataBits:          8,
		StopBits:          1,
		MinimumReadSize:   1,
	}

	port, err := serial.Open(options)
	check(err)

	defer port.Close()

	// Send magic byte to bootloader, to skip 60s wait after startup.
	//
	// "I found the solution.
	//  The serial boot loader is waiting for 60 seconds before jumping to the ZNP application image, as documented in Serial Boot Loader document.
	//  Sending a magic byte 0xef can force the sbl to skip the waiting.
	//  However, the version in Home-1.2.1 changed the magic byte from 0x10 to 0xef, and the document did not update this.
	//  That's why I could not get it work before.
	//  Also the response to the magic byte has been changed from 0x00 to SYS_RESET_IND. again no document update."
	//
	//    - Kyle Zhou in [2]
	//
	// See: [1] https://github.com/Koenkk/zigbee2mqtt/issues/1343
	// See: [2] https://e2e.ti.com/support/wireless-connectivity/zigbee-and-thread/f/158/p/160948/1361000#1361000
	//
	_, err = port.Write([]byte{0xEF})
	check(err)

	_, err = port.Write([]byte{254, 0, 33, 2, 35})
	check(err)

	buffer := make([]byte, 256)
	for {
		n, err := port.Read(buffer)
		check(err)
		fmt.Println(buffer[:n])
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
