package main

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

/* FRAME_SUBSYSTEM_MAC */

/* FRAME_SUBSYSTEM_NWK */

/* FRAME_SUBSYSTEM_AF */

func init() {
	registerCommand(FRAME_TYPE_SREQ, FRAME_SUBSYSTEM_AF, 0x00, AfRegisterRequest{})
	registerCommand(FRAME_TYPE_SRSP, FRAME_SUBSYSTEM_AF, 0x00, AfRegisterResponse{})
}

type AfRegisterRequest struct {
	EndPoint       uint8
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
	DeviceState  uint8
	AssocDevices []uint16
}

/* FRAME_SUBSYSTEM_DEBUG */

/* FRAME_SUBSYSTEM_APP */
