package zigbee

import "fmt"

type DeviceState uint8

const (
	DeviceStateInitializedNotStarted   DeviceState = 0x00 // Initialized - not started automatically
	DeviceStateInitializedNotConnected DeviceState = 0x01 // Initialized - not connected to anything
	DeviceStateDiscovering             DeviceState = 0x02 // Discovering PANs to join
	DeviceStateJoining                 DeviceState = 0x03 // Joining a PAN
	DeviceStateRejoining               DeviceState = 0x04 // Rejoining a PAN, only for end devices
	DeviceStateJoinedNotAuthenticated  DeviceState = 0x05 // Joined but not yet authenticated by trust center
	DeviceStateEndDevice               DeviceState = 0x06 // Started as device after authentication
	DeviceStateRouter                  DeviceState = 0x07 // Device joined, authenticated and is a router
	DeviceStateCoordinatorStarting     DeviceState = 0x08 // Starting as ZigBee Coordinator
	DeviceStateCoordinator             DeviceState = 0x09 // Started as ZigBee Coordinator
	DeviceStateOrphan                  DeviceState = 0x0A // Device has lost information about its parent
)

var DeviceStateNames = []string{
	"InitializedNotStarted",
	"InitializedNotConnected",
	"Discovering",
	"Joining",
	"Rejoining",
	"JoinedNotAuthenticated",
	"EndDevice",
	"Router",
	"CoordinatorStarting",
	"Coordinator",
	"Orphan",
}

func (state DeviceState) String() string {
	if int(state) < len(DeviceStateNames) {
		return fmt.Sprintf("%s(%d)", DeviceStateNames[state], uint8(state))
	}
	return fmt.Sprintf("%d", uint8(state))
}

/* FRAME_SUBSYSTEM_SYS */

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_SYS, 0x02, SysVersionRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_SYS, 0x02, SysVersionResponse{})
}

type SysVersionRequest struct{}

type SysVersionResponse struct {
	TransportRev uint8
	Product      uint8
	MajorRel     uint8
	MinorRel     uint8
	MaintRel     uint8
	Revision     uint32
}

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_SYS, 0x08, SysOsalNvReadRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_SYS, 0x08, SysOsalNvReadResponse{})
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_SYS, 0x09, SysOsalNvWriteRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_SYS, 0x09, SysOsalNvWriteResponse{})
}

type SysOsalNvReadRequest struct {
	ID     uint16
	Offset uint8
}

type SysOsalNvReadResponse struct {
	Status byte
	Value  []byte
}

type SysOsalNvWriteRequest struct {
	ID     uint16
	Offset uint8
	Value  []byte
}

type SysOsalNvWriteResponse struct {
	Status byte
}

func init() {
	registerCommand(FRAME_TYPE_AREQ, FRAME_SUBSYSTEM_SYS, 0x80, SysResetInd{})
}

const (
	ResetReasonPowerUp  = 0x00
	ResetReasonExternal = 0x01
	ResetReasonWatchDog = 0x02
)

type SysResetInd struct {
	Reason       uint8
	TransportRev uint8
	Product      uint8
	MajorRel     uint8
	MinorRel     uint8
	MaintRel     uint8
}

/* FRAME_SUBSYSTEM_MAC */

/* FRAME_SUBSYSTEM_NWK */

/* FRAME_SUBSYSTEM_AF */

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_AF, 0x00, AfRegisterRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_AF, 0x00, AfRegisterResponse{})
}

const (
	LatencyReqNoLatency   = 0x00
	LatencyReqFastBeacons = 0x01
	LatencyReqSlowBeacons = 0x02
)

type AfRegisterRequest struct {
	Endpoint       uint8
	AppProfID      uint16
	AppDeviceID    uint16
	AddDevVer      uint8
	LatencyReq     uint8
	AppInClusters  []uint16
	AppOutClusters []uint16
}

type AfRegisterResponse struct {
	Status byte
}

func init() {
	registerCommand(FRAME_TYPE_AREQ, FRAME_SUBSYSTEM_AF, 0x81, AfIncomingMsg{})
}

type AfIncomingMsg struct {
	GroupID        uint16
	ClusterID      uint16
	SrcAddr        uint16
	SrcEndpoint    uint8
	DstEndpoint    uint8
	WasBroadcast   uint8
	LinkQuality    uint8
	SecureUse      uint8
	TimeStamp      uint32
	TransSeqNumber uint8
	Data           []byte
}

/* FRAME_SUBSYSTEM_ZDO*/

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_ZDO, 0x05, ZdoActiveEPRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_ZDO, 0x05, ZdoActiveEPResponse{})
	registerCommand(FRAME_TYPE_AREQ, FRAME_SUBSYSTEM_ZDO, 0x85, ZdoActiveEP{})
}

type ZdoActiveEPRequest struct {
	DstAddr           uint16
	NWKAddrOfInterest uint16
}

type ZdoActiveEPResponse struct {
	Status byte
}

type ZdoActiveEP struct {
	SrcAddr   uint16
	Status    byte
	NWKAddr   uint16
	ActiveEPs []uint8
}

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_ZDO, 0x36, ZdoMgmtPermitJoinRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_ZDO, 0x36, ZdoMgmtPermitJoinResponse{})
	registerCommand(FRAME_TYPE_AREQ, FRAME_SUBSYSTEM_ZDO, 0xb6, ZdoMgmtPermitJoin{})
}

type ZdoMgmtPermitJoinRequest struct {
	AddrMode       byte
	DstAddr        uint16
	Duration       byte
	TCSignificance byte
}

type ZdoMgmtPermitJoinResponse struct {
	Status byte
}

type ZdoMgmtPermitJoin struct {
	SrcAddr uint16
	Status  byte
}

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_ZDO, 0x40, ZdoStartupFromAppRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_ZDO, 0x40, ZdoStartupFromAppResponse{})
}

type ZdoStartupFromAppRequest struct {
	StartDelay uint16
}

type ZdoStartupFromAppResponse struct {
	Status byte
}

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_ZDO, 0x50, ZdoExtNwkInfoRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_ZDO, 0x50, ZdoExtNwkInfoResponse{})
}

type ZdoExtNwkInfoRequest struct{}

type ZdoExtNwkInfoResponse struct {
	ShortAddress          uint16
	PanID                 uint16
	ParentAddress         uint16
	ExtendedPanID         uint64
	ExtendedParentAddress uint64
	Channel               uint16
}

func init() {
	registerCommand(FRAME_TYPE_AREQ, FRAME_SUBSYSTEM_ZDO, 0xc0, ZdoStateChangeInd{})
}

type ZdoStateChangeInd struct {
	State DeviceState
}

func init() {
	registerCommand(FRAME_TYPE_AREQ, FRAME_SUBSYSTEM_ZDO, 0xc1, ZdoEndDeviceAnnceInd{})
}

const (
	CapabilitiesAlternatePanCoordinator = 1 << 0
	CapabilitiesRouter                  = 1 << 1
	CapabilitiesMainPowered             = 1 << 2
	CapabilitiesReceiverOnWhenIdle      = 1 << 3
	CapabilitiesSecurityCapability      = 1 << 6
)

type ZdoEndDeviceAnnceInd struct {
	SrcAddr      uint16
	NwkAddr      uint16
	IEEEAddr     uint64
	Capabilities byte
}

func init() {
	registerCommand(FRAME_TYPE_AREQ, FRAME_SUBSYSTEM_ZDO, 0xca, ZdoTcDevInd{})
}

// Trust Center Device Indication.
type ZdoTcDevInd struct {
	SrcNwkAddr    uint16
	WasBroadcast  uint64
	ParentNwkAddr uint16
}

func init() {
	registerCommand(FRAME_TYPE_AREQ, FRAME_SUBSYSTEM_ZDO, 0xcb, ZdoPermitJoinInd{})
}

type ZdoPermitJoinInd struct {
	Duration byte
}

/* FRAME_SUBSYSTEM_SAPI */

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_SAPI, 0x04, ZbReadConfigurationRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_SAPI, 0x04, ZBReadConfigurationResponse{})
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_SAPI, 0x05, ZbWriteConfigurationRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_SAPI, 0x05, ZbWriteConfigurationResponse{})
}

type ZbReadConfigurationRequest struct {
	ConfigID uint8
}

type ZBReadConfigurationResponse struct {
	Status   byte
	ConfigID uint8
	Value    []byte
}

type ZbWriteConfigurationRequest struct {
	ConfigID uint8
	Value    []byte
}

type ZbWriteConfigurationResponse struct {
	Status byte
}

/* FRAME_SUBSYSTEM_UTIL */

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_UTIL, 0x00, UtilGetDeviceInfoRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_UTIL, 0x00, UtilGetDeviceInfoResponse{})
}

type UtilGetDeviceInfoRequest struct{}

type UtilGetDeviceInfoResponse struct {
	Status       byte
	IEEEAddr     uint64
	ShortAddr    uint16
	DeviceType   uint8
	DeviceState  DeviceState
	AssocDevices []uint16
}

/* FRAME_SUBSYSTEM_DEBUG */

/* FRAME_SUBSYSTEM_APP */
