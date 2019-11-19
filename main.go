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

type FrameType byte

const (
	FRAME_TYPE_POLL FrameType = 0
	FRAME_TYPE_SREQ FrameType = 1 // synchronous request
	FRAME_TYPE_AREQ FrameType = 2 // asynchronous reqest
	FRAME_TYPE_SRSP FrameType = 3 // synchronous response
)

func (t FrameType) String() string {
	switch t {
	case FRAME_TYPE_POLL:
		return "POLL"
	case FRAME_TYPE_SREQ:
		return "SREQ"
	case FRAME_TYPE_AREQ:
		return "AREQ"
	case FRAME_TYPE_SRSP:
		return "SRSP"
	default:
		return fmt.Sprintf("Reserved(%d)", byte(t))
	}
}

type FrameSubsystem byte

const (
	FRAME_SUBSYSTEM_RPC_ERROR FrameSubsystem = 0
	FRAME_SUBSYSTEM_SYS       FrameSubsystem = 1
	FRAME_SUBSYSTEM_MAC       FrameSubsystem = 2
	FRAME_SUBSYSTEM_NWK       FrameSubsystem = 3
	FRAME_SUBSYSTEM_AF        FrameSubsystem = 4 // application framework
	FRAME_SUBSYSTEM_ZDO       FrameSubsystem = 5 // zigbee device object
	FRAME_SUBSYSTEM_SAPI      FrameSubsystem = 6 // simple api
	FRAME_SUBSYSTEM_UTIL      FrameSubsystem = 7
	FRAME_SUBSYSTEM_DEBUG     FrameSubsystem = 8
	FRAME_SUBSYSTEM_APP       FrameSubsystem = 9
)

func (s FrameSubsystem) String() string {
	switch s {
	case FRAME_SUBSYSTEM_RPC_ERROR:
		return "RPC_ERROR"
	case FRAME_SUBSYSTEM_SYS:
		return "SYS"
	case FRAME_SUBSYSTEM_MAC:
		return "MAC"
	case FRAME_SUBSYSTEM_NWK:
		return "NWK"
	case FRAME_SUBSYSTEM_AF:
		return "AF"
	case FRAME_SUBSYSTEM_ZDO:
		return "ZDO"
	case FRAME_SUBSYSTEM_SAPI:
		return "SAPI"
	case FRAME_SUBSYSTEM_UTIL:
		return "UTIL"
	case FRAME_SUBSYSTEM_DEBUG:
		return "DEBUG"
	case FRAME_SUBSYSTEM_APP:
		return "APP"
	default:
		return fmt.Sprintf("Reserved(%d)", byte(s))
	}
}

type Frame struct {
	Type      FrameType
	Subsystem FrameSubsystem
	ID        byte
	Data      []byte
}

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

	check(writeFrame(port, Frame{FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_SYS, 2, nil}))

	time.Sleep(20 * time.Second)
}

func writeFrame(port io.Writer, frame Frame) error {
	frameLength := 1 + 1 + 2 + len(frame.Data) + 1 // SOF (start of frame) + length + command + data + FCS (frame check sequence)
	var buffer [256]byte
	buffer[0] = SOF
	buffer[1] = byte(len(frame.Data))
	buffer[2] = (byte(frame.Type) << 5) | byte(frame.Subsystem) // cmd0
	buffer[3] = frame.ID                                        // cmd1
	copy(buffer[4:], frame.Data)
	buffer[frameLength-1] = calculateFCS(buffer[1 : frameLength-1]) // exclude SOF and FCS
	_, err := port.Write(buffer[:frameLength])
	return err
}

// @Todo: Ignore errors.Is(err, os.ErrClosed).
func readPort(port io.Reader) {
	r := bufio.NewReaderSize(port, 256)
	for {
		// Find SOF (start of frame).
		b, err := r.ReadByte()
		check(err)
		if b != SOF {
			skipped := 0
			for b != SOF {
				b, err = r.ReadByte()
				check(err)
				skipped++
			}
			log.Printf("skipped %d bytes before frame", skipped)
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

		cmd0, cmd1 := buffer[1], buffer[2]

		frame := Frame{
			Type:      FrameType(cmd0 >> 5),
			Subsystem: FrameSubsystem(cmd0 & 0b00011111),
			ID:        cmd1,
			Data:      buffer[3:],
		}

		fmt.Println("frame: ", frame)
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
