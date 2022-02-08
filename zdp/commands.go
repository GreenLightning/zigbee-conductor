package zdp

import (
	"github.com/GreenLightning/zigbee-conductor/zigbee"
)

type NWKAddrReq struct {
	IEEEAddress zigbee.MACAddress
	RequestType uint8
	StartIndex  uint8
}

type IEEEAddrReq struct {
	NWKAddrOfInterest uint16
	RequestType       uint8
	StartIndex        uint8
}

type ActiveEPReq struct {
	NWKAddrOfInterest uint16
}

type MatchDescReq struct {
	NWKAddrOfInterest uint16
	ProfileID         zigbee.ProfileID
	InClusters        []uint16
	OutClusters       []uint16
}

type DeviceAnnce struct {
	NWKAddr    uint16
	IEEEAddr   zigbee.MACAddress
	Capability uint8
}

type ParentAnnce struct {
	ChildInfo []zigbee.MACAddress
}

type EndDeviceBindReq struct {
	BindingTarget  uint16
	SrcIEEEAddress zigbee.MACAddress
	SrcEndpoint    uint8
	ProfileID      zigbee.ProfileID
	InClusters     []uint16
	OutClusters    []uint16
}

type BindReq struct {
	SrcAddress  zigbee.MACAddress
	SrcEndpoint uint8
	ClusterID   uint16
	DstAddress  zigbee.Address
	DstEndpoint uint8
}

type UnbindReq struct {
	SrcAddress  zigbee.MACAddress
	SrcEndpoint uint8
	ClusterID   uint16
	DstAddress  zigbee.Address
	DstEndpoint uint8
}

type BindRegisterReq struct {
	NodeAddress zigbee.MACAddress
}

type ReplaceDeviceReq struct {
	OldAddress  zigbee.MACAddress
	OldEndpoint uint8
	NewAddress  zigbee.MACAddress
	NewEndpoint uint8
}

type NWKAddrRsp struct {
	Status            Status
	IEEEAddrRemoteDev zigbee.MACAddress
	NWKAddrRemoteDev  uint16
	StartIndex        uint8
	NWKAddrAssocDevs  []uint16
}

type IEEEAddrRsp struct {
	Status            Status
	IEEEAddrRemoteDev zigbee.MACAddress
	NWKAddrRemoteDev  uint16
	StartIndex        uint8
	NWKAddrAssocDevs  []uint16
}

type ActiveEPRsp struct {
	Status            Status
	NWKAddrOfInterest uint16
	ActiveEPs         []uint8
}

type MatchDescRsp struct {
	Status            Status
	NWKAddrOfInterest uint16
	Matches           []uint8
}

type ParentAnnceRsp struct {
	Status    Status
	ChildInfo []zigbee.MACAddress
}

type EndDeviceBindRsp struct {
	Status Status
}

type BindRsp struct {
	Status Status
}

type UnbindRsp struct {
	Status Status
}

type BindRegisterRsp struct {
	Status              Status
	BindingTableEntries uint16
	BindingTableCount   uint16
	BindingTable        []uint8 // ?
}

type ReplaceDeviceRsp struct {
	Status Status
}
