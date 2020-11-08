package zcl

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type FrameType byte

const (
	FrameTypeGlobal FrameType = 0b00
	FrameTypeLocal  FrameType = 0b01
)

func (t FrameType) String() string {
	switch t {
	case FrameTypeGlobal:
		return "Global"
	case FrameTypeLocal:
		return "Local"
	default:
		return fmt.Sprintf("Reserved(%d)", byte(t))
	}
}

type FrameHeader struct {
	Type                    FrameType
	ManufacturerSpecific    bool
	DirectionServerToClient bool
	DisableDefaultResponse  bool

	ManufacturerCode uint16
	TransSeqNumber   uint8
	CommandID        CommandID
}

type Frame struct {
	FrameHeader
	Data []byte
}

func ParseFrame(data []byte) (Frame, error) {
	var frame Frame

	if len(data) < 1 {
		return frame, errors.New("no data")
	}

	control := data[0]
	data = data[1:]
	frame.Type = FrameType(control & 0b00011)
	frame.ManufacturerSpecific = (control & 0b00100) != 0
	frame.DirectionServerToClient = (control & 0b01000) != 0
	frame.DisableDefaultResponse = (control & 0b10000) != 0

	if frame.ManufacturerSpecific {
		if len(data) < 2 {
			return frame, errors.New("not enough data")
		}
		frame.ManufacturerCode = binary.LittleEndian.Uint16(data)
		data = data[2:]
	}

	if len(data) < 2 {
		return frame, errors.New("not enough data")
	}

	frame.TransSeqNumber = data[0]
	frame.CommandID = CommandID(data[1])
	data = data[2:]

	frame.Data = data
	return frame, nil
}

func SerializeFrame(frame Frame) []byte {
	length := 3 + len(frame.Data)
	if frame.ManufacturerSpecific {
		length += 2
	}

	data := make([]byte, 0, length)

	var control byte
	control &= byte(frame.Type) & 0b00011
	if frame.ManufacturerSpecific {
		control &= 0b00100
	}
	if frame.DirectionServerToClient {
		control &= 0b01000
	}
	if frame.DisableDefaultResponse {
		control &= 0b10000
	}
	data = append(data, control)

	if frame.ManufacturerSpecific {
		start := len(data)
		data = append(data, 0, 0)
		binary.LittleEndian.PutUint16(data[start:], frame.ManufacturerCode)
	}

	data = append(data, frame.TransSeqNumber)
	data = append(data, byte(frame.CommandID))

	data = append(data, frame.Data...)
	return data
}
