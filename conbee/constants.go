package conbee

import (
	"fmt"
	"strings"
)

type CommandID byte

const (
	CmdDeviceState         CommandID = 0x07
	CmdChangeNetworkState  CommandID = 0x08
	CmdReadParameter       CommandID = 0x0A
	CmdWriteParameter      CommandID = 0x0B
	CmdDeviceStateChanged  CommandID = 0x0E
	CmdVersion             CommandID = 0x0D
	CmdAPSDataRequest      CommandID = 0x12
	CmdAPSDataConfirm      CommandID = 0x04
	CmdAPSDataIndication   CommandID = 0x17
	CmdMACPollIndication   CommandID = 0x1C
	CmdUpdateNeighbor      CommandID = 0x1D
	CmdMACBeaconIndication CommandID = 0x1F
	CmdUpdateBootloader    CommandID = 0x21
)

func (c CommandID) String() string {
	switch c {
	case CmdDeviceState:
		return "DeviceState"
	case CmdChangeNetworkState:
		return "ChangeNetworkState"
	case CmdReadParameter:
		return "ReadParameter"
	case CmdWriteParameter:
		return "WriteParameter"
	case CmdDeviceStateChanged:
		return "DeviceStateChanged"
	case CmdVersion:
		return "Version"
	case CmdAPSDataRequest:
		return "APSDataRequest"
	case CmdAPSDataConfirm:
		return "APSDataConfirm"
	case CmdAPSDataIndication:
		return "APSDataIndication"
	case CmdMACPollIndication:
		return "MACPollIndication"
	case CmdUpdateNeighbor:
		return "UpdateNeighbor"
	case CmdMACBeaconIndication:
		return "MACBeaconIndication"
	case CmdUpdateBootloader:
		return "UpdateBootloader"
	default:
		return fmt.Sprintf("CommandID(%02x)", byte(c))
	}
}

type Status byte

const (
	StatusSuccess      Status = 0x00
	StatusFailure      Status = 0x01
	StatusBusy         Status = 0x02
	StatusTimeout      Status = 0x03
	StatusUnsupported  Status = 0x04
	StatusError        Status = 0x05
	StatusNoNetwork    Status = 0x06
	StatusInvalidValue Status = 0x07
)

func (s Status) String() string {
	switch s {
	case StatusSuccess:
		return "Success"
	case StatusFailure:
		return "Failure"
	case StatusBusy:
		return "Busy"
	case StatusTimeout:
		return "Timeout"
	case StatusUnsupported:
		return "Unsupported"
	case StatusError:
		return "Error"
	case StatusNoNetwork:
		return "NoNetwork"
	case StatusInvalidValue:
		return "InvalidValue"
	default:
		return fmt.Sprintf("Status(%02x)", byte(s))
	}
}

type DeviceState byte

const (
	DeviceStateNetworkStateMask         DeviceState = 0b0000_0011
	DeviceStateDataConfirmFlag          DeviceState = 0b0000_0100
	DeviceStateDataIndicationFlag       DeviceState = 0b0000_1000
	DeviceStateConfigurationChangedFlag DeviceState = 0b0001_0000
	DeviceStateDataRequestFreeSlotsFlag DeviceState = 0b0010_0000
)

func (s DeviceState) NetworkState() NetworkState {
	return NetworkState(s & DeviceStateNetworkStateMask)
}

func (s DeviceState) String() string {
	var builder strings.Builder

	if s&DeviceStateDataRequestFreeSlotsFlag != 0 {
		if builder.Len() != 0 {
			builder.WriteString("|")
		}
		builder.WriteString("DataRequestFreeSlots")
	}
	if s&DeviceStateConfigurationChangedFlag != 0 {
		if builder.Len() != 0 {
			builder.WriteString("|")
		}
		builder.WriteString("ConfChanged")
	}
	if s&DeviceStateDataIndicationFlag != 0 {
		if builder.Len() != 0 {
			builder.WriteString("|")
		}
		builder.WriteString("DataInd")
	}
	if s&DeviceStateDataConfirmFlag != 0 {
		if builder.Len() != 0 {
			builder.WriteString("|")
		}
		builder.WriteString("DataConfirm")
	}

	if builder.Len() != 0 {
		builder.WriteString("|")
	}
	builder.WriteString(s.NetworkState().String())

	return builder.String()
}

type NetworkState byte

const (
	NetworkStateOffline   NetworkState = 0x00
	NetworkStateJoining   NetworkState = 0x01
	NetworkStateConnected NetworkState = 0x02
	NetworkStateLeaving   NetworkState = 0x03
)

func (s NetworkState) String() string {
	switch s {
	case NetworkStateOffline:
		return "Offline"
	case NetworkStateJoining:
		return "Joining"
	case NetworkStateConnected:
		return "Connected"
	case NetworkStateLeaving:
		return "Leaving"
	default:
		return fmt.Sprintf("NetworkState(%02x)", byte(s))
	}
}

type NetParam byte

const (
	NetParamMACAddress             NetParam = 0x01
	NetParamNWKPANID               NetParam = 0x05
	NetParamNWKAddress             NetParam = 0x07
	NetParamNWKExtendedPANID       NetParam = 0x08
	NetParamAPSDesignedCoordinator NetParam = 0x09
	NetParamChannelMask            NetParam = 0x0A
	NetParamAPSExtendedPANID       NetParam = 0x0B
	NetParamTrustCenterAddress     NetParam = 0x0E
	NetParamSecurityMode           NetParam = 0x10
	NetParamPredefinedNWKPANID     NetParam = 0x15
	NetParamNetworkKey             NetParam = 0x18
	NetParamLinkKey                NetParam = 0x19
	NetParamCurrentChannel         NetParam = 0x1C
	NetParamProtocolVersion        NetParam = 0x22
	NetParamNWKUpdateID            NetParam = 0x24
	NetParamWatchdogTTL            NetParam = 0x26
	NetParamNWKFrameCounter        NetParam = 0x27
)

func (p NetParam) String() string {
	switch p {
	case NetParamMACAddress:
		return "MACAddress"
	case NetParamNWKPANID:
		return "NWKPANID"
	case NetParamNWKAddress:
		return "NWKAddress"
	case NetParamNWKExtendedPANID:
		return "NWKExtendedPANID"
	case NetParamAPSDesignedCoordinator:
		return "APSDesignedCoordinator"
	case NetParamChannelMask:
		return "ChannelMask"
	case NetParamAPSExtendedPANID:
		return "APSExtendedPANID"
	case NetParamTrustCenterAddress:
		return "TrustCenterAddress"
	case NetParamSecurityMode:
		return "SecurityMode"
	case NetParamPredefinedNWKPANID:
		return "PredefinedNWKPANID"
	case NetParamNetworkKey:
		return "NetworkKey"
	case NetParamLinkKey:
		return "LinkKey"
	case NetParamCurrentChannel:
		return "CurrentChannel"
	case NetParamProtocolVersion:
		return "ProtocolVersion"
	case NetParamNWKUpdateID:
		return "NWKUpdateID"
	case NetParamWatchdogTTL:
		return "WatchdogTTL"
	case NetParamNWKFrameCounter:
		return "NWKFrameCounter"
	default:
		return fmt.Sprintf("NetParam(%02x)", byte(p))
	}
}
