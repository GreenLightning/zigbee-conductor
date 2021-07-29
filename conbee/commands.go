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

type UpdateNeighborCommand struct {
	Action byte
	ShortAddress uint16
	MACAddress MACAddress
}

func init() {
	registerParsable(incomingParsables, CmdUpdateNeighbor, new(UpdateNeighborCommand))
	registerParsable(outgoingParsables, CmdUpdateNeighbor, new(UpdateNeighborCommand))
}

func (c *UpdateNeighborCommand) ParsePayload(data []byte) error {
	if len(data) < 13 {
		return ErrInvalidPacket
	}
	payloadLength := uint16(binary.LittleEndian.Uint16(data))
	if int(payloadLength) + 2 != len(data) {
		return ErrInvalidPacket
	}
	c.Action = data[2]
	c.ShortAddress = binary.LittleEndian.Uint16(data[3:5])
	c.MACAddress = MACAddress(binary.LittleEndian.Uint64(data[5:13]))
	return nil
}
