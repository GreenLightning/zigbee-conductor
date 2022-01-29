package znp

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/GreenLightning/zigbee-conductor/pkg/scf"
)

var (
	frameHeaderByCommandType = make(map[reflect.Type]FrameHeader)
	commandTypeByFrameHeader = make(map[FrameHeader]reflect.Type)
)

func registerCommand(frameType FrameType, subsystem FrameSubsystem, id byte, commandPrototype interface{}) {
	commandType := reflect.TypeOf(commandPrototype)
	if err := scf.ValidateType(commandType); err != nil {
		panic(err)
	}

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

func getHeaderForCommand(command interface{}) FrameHeader {
	commandType := reflect.TypeOf(command)
	return getHeaderForCommandType(commandType)
}

func getHeaderForCommandType(commandType reflect.Type) FrameHeader {
	header, ok := frameHeaderByCommandType[commandType]
	if !ok {
		panic(fmt.Sprintf("command %s has never been registered", commandType.Name()))
	}

	return header
}

func buildFrameForCommand(command interface{}) Frame {
	commandType := reflect.TypeOf(command)
	frame := Frame{FrameHeader: getHeaderForCommandType(commandType)}
	frame.Data = scf.Serialize(command)
	if len(frame.Data) > FRAME_MAX_DATA_LENGTH {
		panic(fmt.Sprintf("command too large: %#v", command))
	}
	return frame
}

var (
	ErrCommandUnknownFrameHeader = errors.New("unknown serial frame header")
	ErrCommandInvalidFrame       = errors.New("invalid serial frame")
)

func parseCommandFromFrame(frame Frame) (interface{}, error) {
	commandType, ok := commandTypeByFrameHeader[frame.FrameHeader]
	if !ok {
		return nil, ErrCommandUnknownFrameHeader
	}

	commandValue := reflect.New(commandType).Elem()
	_, err := scf.ParseValue(commandValue, frame.Data)
	if err != nil {
		return nil, ErrCommandInvalidFrame
	}

	return commandValue.Interface(), nil
}
