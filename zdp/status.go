package zdp

import "fmt"

type Status uint8

const (
	StatusSuccess                Status = 0x00
	StatusInvalidRequestType     Status = 0x80
	StatusDeviceNotFound         Status = 0x81
	StatusInvalidEP              Status = 0x82
	StatusNotActive              Status = 0x83
	StatusNotSupported           Status = 0x84
	StatusTimeout                Status = 0x85
	StatusNoMatch                Status = 0x86
	StatusNoEntry                Status = 0x88
	StatusNoDescriptor           Status = 0x89
	StatusInsufficientSpace      Status = 0x8a
	StatusNotPermitted           Status = 0x8b
	StatusTableFull              Status = 0x8c
	StatusNotAuthorized          Status = 0x8d
	StatusDeviceBindingTableFull Status = 0x8e
	StatusInvalidIndex           Status = 0x8f
)

func (status Status) String() string {
	switch status {
	case StatusSuccess:
		return "Success"
	case StatusInvalidRequestType:
		return "InvalidRequestType"
	case StatusDeviceNotFound:
		return "DeviceNotFound"
	case StatusInvalidEP:
		return "InvalidEP"
	case StatusNotActive:
		return "NotActive"
	case StatusNotSupported:
		return "NotSupported"
	case StatusTimeout:
		return "Timeout"
	case StatusNoMatch:
		return "NoMatch"
	case StatusNoEntry:
		return "NoEntry"
	case StatusNoDescriptor:
		return "NoDescriptor"
	case StatusInsufficientSpace:
		return "InsufficientSpace"
	case StatusNotPermitted:
		return "NotPermitted"
	case StatusTableFull:
		return "TableFull"
	case StatusNotAuthorized:
		return "NotAuthorized"
	case StatusDeviceBindingTableFull:
		return "DeviceBindingTableFull"
	case StatusInvalidIndex:
		return "InvalidIndex"
	default:
		return fmt.Sprintf("0x%02x", uint8(status))
	}
}
