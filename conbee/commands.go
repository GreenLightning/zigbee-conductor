package conbee

import (
	"encoding/binary"
	"fmt"
)

type VersionNumber uint32

func (v VersionNumber) String() string {
	return fmt.Sprintf("0x%08x", uint32(v))
}

type MACAddress uint64

func (a MACAddress) String() string {
	return fmt.Sprintf("%016x", uint64(a))
}

// READ FIRMWARE VERSION

type ReadFirmwareVersionRequest struct {
	Reserved uint32
}

func init() {
	registerParsable(outgoingParsables, CmdVersion, new(ReadFirmwareVersionRequest))
}

func (r *ReadFirmwareVersionRequest) ParsePayload(data []byte) error {
	if len(data) < 4 {
		return ErrInvalidPacket
	}
	r.Reserved = binary.LittleEndian.Uint32(data)
	return nil
}

type ReadFirmwareVersionResponse struct {
	Version VersionNumber
}

func init() {
	registerParsable(incomingParsables, CmdVersion, new(ReadFirmwareVersionResponse))
}

func (r *ReadFirmwareVersionResponse) ParsePayload(data []byte) error {
	if len(data) < 4 {
		return ErrInvalidPacket
	}
	r.Version = VersionNumber(binary.LittleEndian.Uint32(data))
	return nil
}

// READ PARAMETER

type ReadParameterRequest struct {
	ParameterID NetParam
}

func init() {
	registerParsable(outgoingParsables, CmdReadParameter, new(ReadParameterRequest))
}

func (r *ReadParameterRequest) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 1)
	if err != nil {
		return err
	}
	r.ParameterID = NetParam(data[0])
	return nil
}

type ReadParameterResponse struct {
	ParameterID NetParam
	Parameter   []byte
}

func init() {
	registerParsable(incomingParsables, CmdReadParameter, new(ReadParameterResponse))
}

func (r *ReadParameterResponse) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 0)
	if err != nil {
		return err
	}
	// @Note: Payload might be empty, if the parameter is not supported.
	if len(data) != 0 {
		r.ParameterID = NetParam(data[0])
		r.Parameter = data[1:]
	}
	return nil
}

// WRITE PARAMETER

type WriteParameterRequest struct {
	ParameterID NetParam
	Parameter   []byte
}

func init() {
	registerParsable(outgoingParsables, CmdWriteParameter, new(WriteParameterRequest))
}

func (r *WriteParameterRequest) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 1)
	if err != nil {
		return err
	}
	r.ParameterID = NetParam(data[0])
	r.Parameter = data[1:]
	return nil
}

type WriteParameterResponse struct {
	ParameterID NetParam
}

func init() {
	registerParsable(incomingParsables, CmdWriteParameter, new(WriteParameterResponse))
}

func (r *WriteParameterResponse) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 1)
	if err != nil {
		return err
	}
	r.ParameterID = NetParam(data[0])
	return nil
}

// DEVICE STATE

type DeviceStateRequest struct{}

func init() {
	registerParsable(outgoingParsables, CmdDeviceState, new(DeviceStateRequest))
}

func (r *DeviceStateRequest) ParsePayload(data []byte) error {
	return nil
}

type DeviceStateResponse struct {
	State DeviceState
}

func init() {
	registerParsable(incomingParsables, CmdDeviceState, new(DeviceStateResponse))
}

func (r *DeviceStateResponse) ParsePayload(data []byte) error {
	if len(data) < 1 {
		return ErrInvalidPacket
	}
	r.State = DeviceState(data[0])
	return nil
}

// UPDATE NEIGHBOR

type UpdateNeighborCommand struct {
	Action       byte
	ShortAddress uint16
	MACAddress   MACAddress
}

func init() {
	registerParsable(incomingParsables, CmdUpdateNeighbor, new(UpdateNeighborCommand))
	registerParsable(outgoingParsables, CmdUpdateNeighbor, new(UpdateNeighborCommand))
}

func (c *UpdateNeighborCommand) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 11)
	if err != nil {
		return err
	}
	c.Action = data[0]
	c.ShortAddress = binary.LittleEndian.Uint16(data[1:3])
	c.MACAddress = MACAddress(binary.LittleEndian.Uint64(data[3:11]))
	return nil
}

// UTILITIES

func extractPayload(data []byte, minimumPayloadLength int) ([]byte, error) {
	if len(data) < 2+minimumPayloadLength {
		return nil, ErrInvalidPacket
	}
	payloadLength := uint16(binary.LittleEndian.Uint16(data))
	if len(data) != 2+int(payloadLength) {
		return nil, ErrInvalidPacket
	}
	return data[2:], nil
}
