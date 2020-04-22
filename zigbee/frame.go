package zigbee

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

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

type FrameID byte

func (id FrameID) String() string {
	return fmt.Sprintf("0x%x", byte(id))
}

type FrameHeader struct {
	Type      FrameType
	Subsystem FrameSubsystem
	ID        FrameID
}

type Frame struct {
	FrameHeader
	Data []byte
}

const FRAME_MAX_DATA_LENGTH = 250

const SOF = 0xFE // start of frame

func writeFrame(writer io.Writer, frame Frame) error {
	if len(frame.Data) > FRAME_MAX_DATA_LENGTH {
		panic("frame too large")
	}
	frameLength := 1 + 1 + 2 + len(frame.Data) + 1 // SOF (start of frame) + length + command + data + FCS (frame check sequence)
	var buffer [256]byte
	buffer[0] = SOF
	buffer[1] = byte(len(frame.Data))
	buffer[2] = (byte(frame.Type) << 5) | byte(frame.Subsystem) // cmd0
	buffer[3] = byte(frame.ID)                                  // cmd1
	copy(buffer[4:], frame.Data)
	buffer[frameLength-1] = calculateFCS(buffer[1 : frameLength-1]) // exclude SOF and FCS
	_, err := writer.Write(buffer[:frameLength])
	return err
}

var (
	ErrInvalidFrame = errors.New("invalid frame")
	ErrGarbage      = errors.New("garbage")
)

func readFrame(r *bufio.Reader) (Frame, error) {
	var frame Frame

	// Find SOF (start of frame).
	b, err := r.ReadByte()
	if err != nil {
		return frame, err
	}
	if b != SOF {
		skipped := 0
		for b != SOF {
			b, err = r.ReadByte()
			if err != nil {
				return frame, err
			}
			skipped++
		}
		r.UnreadByte()
		return frame, fmt.Errorf("found %d bytes of %w", skipped, ErrGarbage)
	}

	length, err := r.ReadByte()
	if err != nil {
		return frame, err
	}

	buffer := make([]byte, 1+2+length) // length + command + data
	buffer[0] = length
	_, err = io.ReadFull(r, buffer[1:])
	if err != nil {
		return frame, err
	}

	fcs, err := r.ReadByte()
	if err != nil {
		return frame, err
	}

	if fcs != calculateFCS(buffer) {
		return frame, ErrInvalidFrame
	}

	cmd0, cmd1 := buffer[1], buffer[2]

	frame.FrameHeader = FrameHeader{
		Type:      FrameType(cmd0 >> 5),
		Subsystem: FrameSubsystem(cmd0 & 0b00011111),
		ID:        FrameID(cmd1),
	}
	frame.Data = buffer[3:]
	return frame, nil
}

// calculateFCS returns the frame check sequence for a general format frame.
func calculateFCS(buffer []byte) (fcs byte) {
	for _, value := range buffer {
		fcs ^= value
	}
	return
}
