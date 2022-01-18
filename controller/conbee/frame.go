package conbee

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type ParsableCommand interface {
	CommandID() CommandID
	ParsePayload(data []byte) error
}

type SerializableCommand interface {
	CommandID() CommandID
	SerializePayload(buffer *bytes.Buffer) error
}

type ParsableMap map[CommandID]reflect.Type

var (
	incomingParsables = make(ParsableMap)
	outgoingParsables = make(ParsableMap)
)

func registerParsable(pm ParsableMap, prototype ParsableCommand) {
	id := prototype.CommandID()
	commandType := reflect.TypeOf(prototype)
	if commandType.Kind() != reflect.Ptr || commandType.Elem().Kind() != reflect.Struct {
		panic("command must be a pointer to a struct")
	}
	if old, ok := pm[id]; ok {
		panic(fmt.Sprintf("parsable for %v already registered: old=%s, new=%s", id, old.Name(), commandType.Elem().Name()))
	}
	pm[id] = commandType.Elem()
}

func makeParsable(incoming bool, id CommandID) ParsableCommand {
	pm := incomingParsables
	if !incoming {
		pm = outgoingParsables
	}
	if parsableType, ok := pm[id]; ok {
		return reflect.New(parsableType).Interface().(ParsableCommand)
	}
	return nil
}

// One frame of the Conbee serial protocol.
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

	if value := makeParsable(incoming, frame.CommandID); value != nil {
		frame.Command = value
		err = value.ParsePayload(payload)
		if err != nil {
			return
		}
	} else {
		frame.Command = payload
	}

	return
}

// SerializeFrame converts a frame into a slice of bytes.
//
// The frame length is computed automatically and the corresponding field from
// the frame object is ignored.
//
// The command can be a SerializableCommand or a slice of bytes. If the command
// is a SerializableCommand the command ID is taken from the command.
func SerializeFrame(frame Frame) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.Grow(32)

	if command, ok := frame.Command.(SerializableCommand); ok {
		frame.CommandID = command.CommandID()
	}

	buffer.WriteByte(byte(frame.CommandID))
	buffer.WriteByte(byte(frame.SequenceNumber))
	buffer.WriteByte(byte(frame.Status))

	// Placeholder for frame length.
	buffer.WriteByte(0)
	buffer.WriteByte(0)

	if command, ok := frame.Command.(SerializableCommand); ok {
		err := command.SerializePayload(&buffer)
		if err != nil {
			return nil, err
		}
	} else if command, ok := frame.Command.([]byte); ok {
		buffer.Write(command)
	} else {
		return nil, fmt.Errorf("cannot serialize command: %T", frame.Command)
	}

	frameLength := uint16(buffer.Len())
	binary.LittleEndian.PutUint16(buffer.Bytes()[3:5], frameLength)

	crc := computeCRC(buffer.Bytes())
	WriteUint16(&buffer, crc)

	return buffer.Bytes(), nil
}

func BeginPayload(buffer *bytes.Buffer) int {
	position := buffer.Len()
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	return position
}

func EndPayload(buffer *bytes.Buffer, pos int) {
	length := buffer.Len() - pos - 2
	binary.LittleEndian.PutUint16(buffer.Bytes()[pos:], uint16(length))
}

func WriteUint16(buffer *bytes.Buffer, value uint16) {
	var data [2]byte
	binary.LittleEndian.PutUint16(data[:], value)
	buffer.Write(data[:])
}

func computeCRC(data []byte) (crc uint16) {
	for _, value := range data {
		crc += uint16(value)
	}
	return -crc
}
