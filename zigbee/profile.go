package zigbee

import "fmt"

type ProfileID uint16

const (
	// ZigBee Device Profile
	ProfileDevice                       ProfileID = 0x0000
	ProfileIndustrialPlantMonitoring    ProfileID = 0x0101
	ProfileHomeAutomation               ProfileID = 0x0104
	ProfileCommercialBuildingAutomation ProfileID = 0x0105
	ProfileTelecomApplications          ProfileID = 0x0107
	ProfilePersonalHomeAndHospitalCare  ProfileID = 0x0108
	ProfileAdvancedMeteringInitialtive  ProfileID = 0x0109
)

func (id ProfileID) String() string {
	switch id {
	case ProfileDevice:
		return "ZDP"
	case ProfileIndustrialPlantMonitoring:
		return "IPM"
	case ProfileHomeAutomation:
		return "HA"
	case ProfileCommercialBuildingAutomation:
		return "CBA"
	case ProfileTelecomApplications:
		return "TA"
	case ProfilePersonalHomeAndHospitalCare:
		return "PHHC"
	case ProfileAdvancedMeteringInitialtive:
		return "AMI"
	default:
		return fmt.Sprintf("0x%04x", uint16(id))
	}
}
