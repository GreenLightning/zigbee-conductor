package conbee

import (
	"bytes"
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

type AddressMode byte

const (
	AddressModeNone     AddressMode = 0x00
	AddressModeGroup    AddressMode = 0x01
	AddressModeNWK      AddressMode = 0x02
	AddressModeIEEE     AddressMode = 0x03
	AddressModeCombined AddressMode = 0x04
)

func (m AddressMode) String() string {
	switch m {
	case AddressModeNone:
		return "None"
	case AddressModeGroup:
		return "Group"
	case AddressModeNWK:
		return "NWK"
	case AddressModeIEEE:
		return "IEEE"
	case AddressModeCombined:
		return "Combined"
	default:
		return fmt.Sprintf("AddressMode(%d)", byte(m))
	}
}

type Address struct {
	Mode     AddressMode
	Short    uint16
	Extended MACAddress
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

func (r *ReadFirmwareVersionRequest) SerializePayload(buffer *bytes.Buffer) error {
	var data [4]byte
	binary.LittleEndian.PutUint32(data[:], r.Reserved)
	buffer.Write(data[:])
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

func (r *ReadFirmwareVersionResponse) SerializePayload(buffer *bytes.Buffer) error {
	var data [4]byte
	binary.LittleEndian.PutUint32(data[:], uint32(r.Version))
	buffer.Write(data[:])
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

func (r *ReadParameterRequest) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)
	buffer.WriteByte(byte(r.ParameterID))
	EndPayload(buffer, payload)
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

func (r *ReadParameterResponse) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)
	buffer.WriteByte(byte(r.ParameterID))
	buffer.Write(r.Parameter)
	EndPayload(buffer, payload)
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

func (r *WriteParameterRequest) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)
	buffer.WriteByte(byte(r.ParameterID))
	buffer.Write(r.Parameter)
	EndPayload(buffer, payload)
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

func (r *WriteParameterResponse) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)
	buffer.WriteByte(byte(r.ParameterID))
	EndPayload(buffer, payload)
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

func (r *DeviceStateRequest) SerializePayload(buffer *bytes.Buffer) error {
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(0)
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

func (r *DeviceStateResponse) SerializePayload(buffer *bytes.Buffer) error {
	buffer.WriteByte(byte(r.State))
	buffer.WriteByte(0)
	return nil
}

// RECEIVING DATA

type ReceivedDataNotification struct {
	State DeviceState
}

func init() {
	registerParsable(incomingParsables, CmdDeviceStateChanged, new(ReceivedDataNotification))
}

func (r *ReceivedDataNotification) ParsePayload(data []byte) error {
	if len(data) < 1 {
		return ErrInvalidPacket
	}
	r.State = DeviceState(data[0])
	return nil
}

func (r *ReceivedDataNotification) SerializePayload(buffer *bytes.Buffer) error {
	buffer.WriteByte(byte(r.State))
	buffer.WriteByte(0)
	return nil
}

type ReadReceivedDataRequest struct {
	Flags byte
}

func init() {
	registerParsable(outgoingParsables, CmdAPSDataIndication, new(ReadReceivedDataRequest))
}

func (r *ReadReceivedDataRequest) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 0)
	if err != nil {
		return err
	}
	if len(data) != 0 {
		r.Flags = data[0]
	}
	return nil
}

func (r *ReadReceivedDataRequest) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)
	buffer.WriteByte(byte(r.Flags))
	EndPayload(buffer, payload)
	return nil
}

type ReadReceivedDataResponse struct {
	State               DeviceState
	Destination         Address
	DestinationEndpoint byte
	Source              Address
	SourceEndpoint      byte
	ProfileID           uint16
	ClusterID           uint16

	// Payload is the APS frame payload.
	Payload []byte

	// LQI is the Link Quality Indication.
	LQI byte

	// RSSI is the Received Signal Strength Indication.
	RSSI int8
}

func init() {
	registerParsable(incomingParsables, CmdAPSDataIndication, new(ReadReceivedDataResponse))
}

func (r *ReadReceivedDataResponse) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 1)
	if err != nil {
		return err
	}

	r.State = DeviceState(data[0])
	data = data[1:]

	r.Destination, data, err = extractAddress(data)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return ErrInvalidPacket
	}
	r.DestinationEndpoint = data[0]
	data = data[1:]

	r.Source, data, err = extractAddress(data)
	if err != nil {
		return err
	}

	if len(data) < 7 {
		return ErrInvalidPacket
	}
	r.SourceEndpoint = data[0]
	r.ProfileID = binary.LittleEndian.Uint16(data[1:3])
	r.ClusterID = binary.LittleEndian.Uint16(data[3:5])
	payloadLength := int(binary.LittleEndian.Uint16(data[5:7]))
	data = data[7:]

	if len(data) != payloadLength+8 {
		return ErrInvalidPacket
	}
	r.Payload = data[:payloadLength]
	data = data[payloadLength:]

	// Other fields in this section are reserved and shall be ignored.
	r.LQI = data[2]
	r.RSSI = int8(data[7])

	return nil
}

func (r *ReadReceivedDataResponse) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)

	buffer.WriteByte(byte(r.State))
	if err := writeAddress(buffer, r.Destination); err != nil {
		return err
	}
	buffer.WriteByte(r.DestinationEndpoint)
	if err := writeAddress(buffer, r.Source); err != nil {
		return err
	}
	buffer.WriteByte(r.SourceEndpoint)

	WriteUint16(buffer, r.ProfileID)
	WriteUint16(buffer, r.ClusterID)

	apsPayload := BeginPayload(buffer)
	buffer.Write(r.Payload)
	EndPayload(buffer, apsPayload)

	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(byte(r.LQI))
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(byte(r.RSSI))

	EndPayload(buffer, payload)
	return nil
}

type MACPollIndication struct {
	Source Address

	// LQI is the Link Quality Indication.
	LQI byte

	// RSSI is the Received Signal Strength Indication.
	RSSI int8
}

func init() {
	registerParsable(incomingParsables, CmdMACPollIndication, new(MACPollIndication))
}

func (r *MACPollIndication) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 0)
	if err != nil {
		return err
	}

	r.Source, data, err = extractAddress(data)
	if err != nil {
		return err
	}

	if len(data) < 2 {
		return ErrInvalidPacket
	}
	r.LQI = data[0]
	r.RSSI = int8(data[1])
	data = data[2:]
	return nil
}

func (r *MACPollIndication) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)
	if err := writeAddress(buffer, r.Source); err != nil {
		return err
	}
	buffer.WriteByte(byte(r.LQI))
	buffer.WriteByte(byte(r.RSSI))
	EndPayload(buffer, payload)
	return nil
}

// SEND DATA

type EnqueueSendDataRequest struct {
	RequestID byte
	// Flags is not documented and should be set to 0.
	// When reading and Flags is not zero, the destination address is not valid.
	Flags               byte
	Destination         Address
	DestinationEndpoint byte
	ProfileID           uint16
	ClusterID           uint16
	SourceEndpoint      byte
	Payload             []byte
	TxOptions           byte
	Radius              byte
}

func init() {
	registerParsable(outgoingParsables, CmdAPSDataRequest, new(EnqueueSendDataRequest))
}

func (r *EnqueueSendDataRequest) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 2)
	if err != nil {
		return err
	}

	r.RequestID = data[0]
	r.Flags = data[1]
	data = data[2:]

	if r.Flags != 0 {
		if len(data) < 5 {
			return ErrInvalidPacket
		}
		data = data[5:]
	} else {
		r.Destination, data, err = extractAddress(data)
		if err != nil {
			return err
		}
	}

	if r.Destination.Mode != AddressModeGroup {
		if len(data) == 0 {
			return ErrInvalidPacket
		}
		r.DestinationEndpoint = data[0]
		data = data[1:]
	}

	if len(data) < 7 {
		return ErrInvalidPacket
	}
	r.ProfileID = binary.LittleEndian.Uint16(data[0:2])
	r.ClusterID = binary.LittleEndian.Uint16(data[2:4])
	r.SourceEndpoint = data[4]
	payloadLength := int(binary.LittleEndian.Uint16(data[5:7]))
	data = data[7:]

	if len(data) != payloadLength+2 {
		return ErrInvalidPacket
	}
	r.Payload = data[:payloadLength]
	data = data[payloadLength:]
	r.TxOptions = data[0]
	r.Radius = data[1]
	return nil
}

func (r *EnqueueSendDataRequest) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)

	buffer.WriteByte(r.RequestID)
	buffer.WriteByte(r.Flags)

	if err := writeAddress(buffer, r.Destination); err != nil {
		return err
	}

	if r.Destination.Mode != AddressModeGroup {
		buffer.WriteByte(r.DestinationEndpoint)
	}

	WriteUint16(buffer, r.ProfileID)
	WriteUint16(buffer, r.ClusterID)
	buffer.WriteByte(r.SourceEndpoint)

	apsPayload := BeginPayload(buffer)
	buffer.Write(r.Payload)
	EndPayload(buffer, apsPayload)

	buffer.WriteByte(r.TxOptions)
	buffer.WriteByte(r.Radius)

	EndPayload(buffer, payload)
	return nil
}

type EnqueueSendDataResponse struct {
	State     DeviceState
	RequestID byte
}

func init() {
	registerParsable(incomingParsables, CmdAPSDataRequest, new(EnqueueSendDataResponse))
}

func (r *EnqueueSendDataResponse) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 2)
	if err != nil {
		return err
	}
	r.State = DeviceState(data[0])
	r.RequestID = data[1]
	return nil
}

func (r *EnqueueSendDataResponse) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)
	buffer.WriteByte(byte(r.State))
	buffer.WriteByte(byte(r.RequestID))
	EndPayload(buffer, payload)
	return nil
}

type QuerySendDataRequest struct{}

func init() {
	registerParsable(outgoingParsables, CmdAPSDataConfirm, new(QuerySendDataRequest))
}

func (r *QuerySendDataRequest) ParsePayload(data []byte) error {
	return nil
}

func (r *QuerySendDataRequest) SerializePayload(buffer *bytes.Buffer) error {
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	return nil
}

type QuerySendDataResponse struct {
	State               DeviceState
	RequestID           byte
	Destination         Address
	DestinationEndpoint byte
	SourceEndpoint      byte
	ConfirmStatus       byte
}

func init() {
	registerParsable(incomingParsables, CmdAPSDataConfirm, new(QuerySendDataResponse))
}

func (r *QuerySendDataResponse) ParsePayload(data []byte) error {
	data, err := extractPayload(data, 1)
	if err != nil {
		return err
	}
	r.State = DeviceState(data[0])
	data = data[1:]

	if len(data) == 0 {
		return nil
	}

	r.RequestID = data[0]
	data = data[1:]

	r.Destination, data, err = extractAddress(data)
	if err != nil {
		return err
	}

	if r.Destination.Mode != AddressModeGroup {
		if len(data) == 0 {
			return ErrInvalidPacket
		}
		r.DestinationEndpoint = data[0]
		data = data[1:]
	}

	if len(data) < 2 {
		return ErrInvalidPacket
	}

	r.SourceEndpoint = data[0]
	r.ConfirmStatus = data[1]

	return nil
}

func (r *QuerySendDataResponse) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)

	buffer.WriteByte(byte(r.State))
	buffer.WriteByte(byte(r.RequestID))

	if err := writeAddress(buffer, r.Destination); err != nil {
		return err
	}

	if r.Destination.Mode != AddressModeGroup {
		buffer.WriteByte(r.DestinationEndpoint)
	}

	buffer.WriteByte(r.SourceEndpoint)
	buffer.WriteByte(r.ConfirmStatus)
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(0)

	EndPayload(buffer, payload)
	return nil
}

// UPDATE NEIGHBOR

// UpdateNeighborCommand is experimental and not officially documented.
// See: https://github.com/dresden-elektronik/deconz-rest-plugin/issues/665#issuecomment-401341502
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

func (c *UpdateNeighborCommand) SerializePayload(buffer *bytes.Buffer) error {
	payload := BeginPayload(buffer)

	var data [12]byte
	data[0] = c.Action
	binary.LittleEndian.PutUint16(data[1:3], c.ShortAddress)
	binary.LittleEndian.PutUint64(data[3:11], uint64(c.MACAddress))
	data[11] = 0x80
	buffer.Write(data[:])

	EndPayload(buffer, payload)
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

func extractAddress(data []byte) (Address, []byte, error) {
	if len(data) == 0 {
		return Address{}, data, ErrInvalidPacket
	}
	var addr Address
	addr.Mode = AddressMode(data[0])
	switch addr.Mode {
	case AddressModeGroup, AddressModeNWK:
		if len(data) < 3 {
			return Address{}, data, ErrInvalidPacket
		}
		addr.Short = binary.LittleEndian.Uint16(data[1:])
		return addr, data[3:], nil

	case AddressModeIEEE:
		if len(data) < 9 {
			return Address{}, data, ErrInvalidPacket
		}
		addr.Extended = MACAddress(binary.LittleEndian.Uint64(data[1:]))
		return addr, data[9:], nil

	case AddressModeCombined:
		if len(data) < 11 {
			return Address{}, data, ErrInvalidPacket
		}
		addr.Short = binary.LittleEndian.Uint16(data[1:])
		addr.Extended = MACAddress(binary.LittleEndian.Uint64(data[3:]))
		return addr, data[11:], nil

	default:
		return Address{}, data, ErrInvalidPacket
	}
}

func writeAddress(buffer *bytes.Buffer, addr Address) error {
	var data [11]byte
	data[0] = byte(addr.Mode)

	switch addr.Mode {
	case AddressModeGroup, AddressModeNWK:
		binary.LittleEndian.PutUint16(data[1:], addr.Short)
		buffer.Write(data[:3])
		return nil

	case AddressModeIEEE:
		binary.LittleEndian.PutUint64(data[1:], uint64(addr.Extended))
		buffer.Write(data[:9])
		return nil

	case AddressModeCombined:
		binary.LittleEndian.PutUint16(data[1:], addr.Short)
		binary.LittleEndian.PutUint64(data[3:], uint64(addr.Extended))
		buffer.Write(data[:11])
		return nil

	default:
		return fmt.Errorf("invalid address mode: %v", addr.Mode)
	}
}
