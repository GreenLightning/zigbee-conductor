package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

const SOF = 0xFE // start of frame

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

	go readPort(port)

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

	_, err = port.Write([]byte{SOF, 0, 33, 2, 35})
	check(err)

	time.Sleep(20 * time.Second)
}

// @Todo: Ignore errors.Is(err, os.ErrClosed).
func readPort(port io.Reader) {
	r := bufio.NewReaderSize(port, 256)
	for {
		// Find SOF (start of frame).
		{
			skipped := 0
			for {
				b, err := r.ReadByte()
				check(err)
				if b == SOF {
					break
				}
				skipped++
			}
			if skipped != 0 {
				log.Printf("skipped %d bytes before frame", skipped)
			}
		}

		length, err := r.ReadByte()
		check(err)

		buffer := make([]byte, 1+2+length) // length + command + data
		buffer[0] = length
		_, err = io.ReadFull(r, buffer[1:])
		check(err)

		fcs, err := r.ReadByte()
		check(err)

		if fcs != calculateFCS(buffer) {
			log.Println("skipping invalid frame")
			continue
		}

		// cmd0, cmd1 := buffer[1], buffer[2]
		fmt.Println("frame: ", buffer)
	}
}

// calculateFCS returns the frame check sequence for a general format frame.
func calculateFCS(buffer []byte) (fcs byte) {
	for _, value := range buffer {
		fcs ^= value
	}
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
