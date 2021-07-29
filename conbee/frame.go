package conbee

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type ParsableCommand interface {
	ParsePayload(data []byte) error
}

type ParsableMap map[CommandID]reflect.Type

var (
	incomingParsables = make(ParsableMap)
	outgoingParsables = make(ParsableMap)
)

func registerParsable(pm ParsableMap, id CommandID, prototype ParsableCommand) {
	commandType := reflect.TypeOf(prototype)
	if commandType.Kind() != reflect.Ptr || commandType.Elem().Kind() != reflect.Struct {
		panic("command must be a pointer to a struct")
	}
	if old, ok := pm[id]; ok {
		panic(fmt.Sprintf("parsable for %v already registered: old=%s, new=%s", id, old.Name(), commandType.Elem().Name()))
	}
	pm[id] = commandType.Elem()
}

func makeParsable(incoming bool, id CommandID) (reflect.Value, bool) {
	pm := incomingParsables
	if !incoming {
		pm = outgoingParsables
	}
	if parsableType, ok := pm[id]; ok {
		return reflect.New(parsableType), true
	}
	return reflect.Value{}, false
}

type Frame struct {
	CommandID      CommandID
	SequenceNumber uint8
	Status         Status
	FrameLength    uint16
	Command        interface{}
}

var ErrInvalidPacket = errors.New("invalid packet")

func ParseFrame(data []byte, incoming bool) (frame Frame, err error) {
	if len(data) < 2 {
		err = ErrInvalidPacket
		return
	}

	crc := binary.LittleEndian.Uint16(data[len(data)-2:])
	data = data[:len(data)-2]

	if crc != computeCRC(data) {
		err = ErrInvalidPacket
		return
	}

	if len(data) < 5 {
		err = ErrInvalidPacket
		return
	}

	frame.CommandID = CommandID(data[0])
	frame.SequenceNumber = data[1]
	frame.Status = Status(data[2])
	frame.FrameLength = binary.LittleEndian.Uint16(data[3:5])
	payload := data[5:]

	if int(frame.FrameLength) != len(data) {
		err = ErrInvalidPacket
		return
	}

	if value, ok := makeParsable(incoming, frame.CommandID); ok {
		err = value.Interface().(ParsableCommand).ParsePayload(payload)
		frame.Command = value.Elem().Interface()
		if err != nil {
			return
		}
	} else {
		frame.Command = payload
	}

	return
}

func computeCRC(data []byte) (crc uint16) {
	for _, value := range data {
		crc += uint16(value)
	}
	return -crc
}
