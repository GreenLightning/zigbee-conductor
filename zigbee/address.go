package zigbee

import "fmt"

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
