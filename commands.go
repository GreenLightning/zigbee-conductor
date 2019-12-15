package main

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
}

type SysOsalNvReadRequest struct {
	ID     uint16
	Offset uint8
}

type SysOsalNvReadResponse struct {
	Status byte
	Value  []byte
}
