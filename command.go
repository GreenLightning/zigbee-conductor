package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var (
	frameHeaderByCommandType = make(map[reflect.Type]FrameHeader)
	commandTypeByFrameHeader = make(map[FrameHeader]reflect.Type)
)

func registerCommand(frameType FrameType, subsystem FrameSubsystem, id byte, commandPrototype interface{}) {
	commandType := reflect.TypeOf(commandPrototype)
	validateCommandType(commandType)

	header := FrameHeader{
		Type:      frameType,
		Subsystem: subsystem,
		ID:        FrameID(id),
	}

	if _, ok := frameHeaderByCommandType[commandType]; ok {
		panic(fmt.Sprintf("command %s already registered", commandType.Name()))
	}
	if _, ok := commandTypeByFrameHeader[header]; ok {
		panic(fmt.Sprintf("command %v already registered", header))
	}
	frameHeaderByCommandType[commandType] = header
	commandTypeByFrameHeader[header] = commandType
}

func validateCommandType(commandType reflect.Type) {
	if commandType.Kind() != reflect.Struct {
		panic("command must be a struct")
	}

	for f := 0; f < commandType.NumField(); f++ {
		field := commandType.Field(f)
		if !validCommandFieldType(field.Type) {
			panic(fmt.Sprintf("command field %s.%s has invalid type", commandType.Name(), field.Name))
		}
	}
}

func validCommandFieldType(fieldType reflect.Type) bool {
	validKinds := []reflect.Kind{
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}
	for _, kind := range validKinds {
		if fieldType.Kind() == kind {
			return true
		}
	}
	if fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Uint8 {
		return true
	}
	return false
}

func buildFrameForCommand(command interface{}) Frame {
	commandType := reflect.TypeOf(command)
	header, ok := frameHeaderByCommandType[commandType]
	if !ok {
		panic(fmt.Sprintf("command %s has never been registered", commandType.Name()))
	}

	frame := Frame{FrameHeader: header}

	commandValue := reflect.ValueOf(command)
	for f := 0; f < commandValue.NumField(); f++ {
		field := commandValue.Field(f)
		fieldKind := field.Kind()

		if fieldKind == reflect.Slice {
			data := field.Bytes()
			frame.Data = append(frame.Data, byte(len(data)))
			frame.Data = append(frame.Data, data...)
			continue
		}

		var value uint64
		switch fieldKind {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = uint64(field.Int())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = field.Uint()
		default:
			panic(fmt.Sprintf("serialization for command field not implemented: kind=%v type=%v field=%s.%s", fieldKind, field.Type(), commandType.Name(), commandType.Field(f).Name))
		}

		switch fieldKind {
		case reflect.Int8, reflect.Uint8:
			frame.Data = append(frame.Data, byte(value))

		case reflect.Int16, reflect.Uint16:
			start := len(frame.Data)
			frame.Data = append(frame.Data, 0, 0)
			binary.LittleEndian.PutUint16(frame.Data[start:], uint16(value))

		case reflect.Int32, reflect.Uint32:
			start := len(frame.Data)
			frame.Data = append(frame.Data, 0, 0, 0, 0)
			binary.LittleEndian.PutUint32(frame.Data[start:], uint32(value))

		case reflect.Int64, reflect.Uint64:
			start := len(frame.Data)
			frame.Data = append(frame.Data, 0, 0, 0, 0, 0, 0, 0, 0)
			binary.LittleEndian.PutUint64(frame.Data[start:], uint64(value))
		}
	}

	if len(frame.Data) > FRAME_MAX_DATA_LENGTH {
		panic(fmt.Sprintf("command too large: %#v", command))
	}

	return frame
}

func writeCommand(writer io.Writer, command interface{}) error {
	frame := buildFrameForCommand(command)
	return writeFrame(writer, frame)
}

var (
	ErrCommandUnknownFrameHeader = errors.New("unknown frame header")
	ErrCommandInvalidFrame       = errors.New("invalid frame")
)

func parseCommandFromFrame(frame Frame) (interface{}, error) {
	commandType, ok := commandTypeByFrameHeader[frame.FrameHeader]
	if !ok {
		return nil, ErrCommandUnknownFrameHeader
	}

	data := frame.Data

	commandValue := reflect.New(commandType).Elem()
	for f := 0; f < commandValue.NumField(); f++ {
		field := commandValue.Field(f)
		fieldKind := field.Kind()

		if fieldKind == reflect.Slice {
			if len(data) == 0 {
				return nil, ErrCommandInvalidFrame
			}

			length := int(data[0])
			data = data[1:]

			if len(data) < length {
				return nil, ErrCommandInvalidFrame
			}

			fieldData := data[:length]
			data = data[length:]

			field.SetBytes(fieldData)
			continue
		}

		length := int(field.Type().Size())
		if len(data) < length {
			return nil, ErrCommandInvalidFrame
		}

		var value uint64
		switch fieldKind {
		case reflect.Int8, reflect.Uint8:
			value = uint64(data[0])

		case reflect.Int16, reflect.Uint16:
			value = uint64(binary.LittleEndian.Uint16(data))

		case reflect.Int32, reflect.Uint32:
			value = uint64(binary.LittleEndian.Uint32(data))

		case reflect.Int64, reflect.Uint64:
			value = uint64(binary.LittleEndian.Uint64(data))

		default:
			panic(fmt.Sprintf("deserialization for command field not implemented: kind=%v type=%v field=%s.%s", fieldKind, field.Type(), commandType.Name(), commandType.Field(f).Name))
		}

		data = data[length:]

		switch fieldKind {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(value))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(value)
		}
	}

	return commandValue.Interface(), nil
}
