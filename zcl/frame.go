package zcl

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type FrameType byte

const (
	FRAME_TYPE_GLOBAL FrameType = 0b00
	FRAME_TYPE_LOCAL  FrameType = 0b01
)

func (t FrameType) String() string {
	switch t {
	case FRAME_TYPE_GLOBAL:
		return "Global"
	case FRAME_TYPE_LOCAL:
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

	ManufacturerCode          uint16
	TransactionSequenceNumber uint8
	CommandIdentifier         uint8
}

type Frame struct {
	ClusterID ClusterID
	FrameHeader
	Data []byte
}

func ParseFrame(clusterID uint16, data []byte) (Frame, error) {
	var frame Frame

	frame.ClusterID = ClusterID(clusterID)

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

	frame.TransactionSequenceNumber = data[0]
	frame.CommandIdentifier = data[1]
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

	data = append(data, frame.TransactionSequenceNumber)
	data = append(data, frame.CommandIdentifier)

	data = append(data, frame.Data...)
	return data
}
