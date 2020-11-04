package zcl

import "fmt"

type ClusterID uint16

const (
	CLUSTER_ID_GENERAL_BASIC                ClusterID = 0x0000
	CLUSTER_ID_GENERAL_POWER_CONFIG         ClusterID = 0x0001
	CLUSTER_ID_GENERAL_DEVICE_TEMP_CONFIG   ClusterID = 0x0002
	CLUSTER_ID_GENERAL_IDENTIFY             ClusterID = 0x0003
	CLUSTER_ID_GENERAL_GROUPS               ClusterID = 0x0004
	CLUSTER_ID_GENERAL_SCENES               ClusterID = 0x0005
	CLUSTER_ID_GENERAL_ON_OFF               ClusterID = 0x0006
	CLUSTER_ID_GENERAL_ON_OFF_SWITCH_CONFIG ClusterID = 0x0007
	CLUSTER_ID_GENERAL_LEVEL_CONTROL        ClusterID = 0x0008
	CLUSTER_ID_GENERAL_ALARMS               ClusterID = 0x0009
	CLUSTER_ID_GENERAL_TIME                 ClusterID = 0x000a
	CLUSTER_ID_GENERAL_LOCATION             ClusterID = 0x000b
	CLUSTER_ID_GENERAL_DIAGNOSTICS          ClusterID = 0x0b05
	CLUSTER_ID_GENERAL_POLL_CONTROL         ClusterID = 0x0020
	CLUSTER_ID_GENERAL_POWER_PROFILE        ClusterID = 0x001a
	CLUSTER_ID_GENERAL_METER_IDENTIFICATION ClusterID = 0x0b01

	CLUSTER_ID_GENERAL_ANALOG_INPUT_BASIC      ClusterID = 0x000c
	CLUSTER_ID_GENERAL_ANALOG_OUTPUT_BASIC     ClusterID = 0x000d
	CLUSTER_ID_GENERAL_ANALOG_VALUE_BASIC      ClusterID = 0x000e
	CLUSTER_ID_GENERAL_BINARY_INPUT_BASIC      ClusterID = 0x000f
	CLUSTER_ID_GENERAL_BINARY_OUTPUT_BASIC     ClusterID = 0x0010
	CLUSTER_ID_GENERAL_BINARY_VALUE_BASIC      ClusterID = 0x0011
	CLUSTER_ID_GENERAL_MULTISTATE_INPUT_BASIC  ClusterID = 0x0012
	CLUSTER_ID_GENERAL_MULTISTATE_OUTPUT_BASIC ClusterID = 0x0013
	CLUSTER_ID_GENERAL_MULTISTATE_VALUE_BASIC  ClusterID = 0x0014

	CLUSTER_ID_MS_ILLUMINANCE_MEASUREMENT          ClusterID = 0x0400
	CLUSTER_ID_MS_ILLUMINANCE_LEVEL_SENSING_CONFIG ClusterID = 0x0401
	CLUSTER_ID_MS_TEMPERATURE_MEASUREMENT          ClusterID = 0x0402
	CLUSTER_ID_MS_PRESSURE_MEASUREMENT             ClusterID = 0x0403
	CLUSTER_ID_MS_FLOW_MEASUREMENT                 ClusterID = 0x0404
	CLUSTER_ID_MS_RELATIVE_HUMIDITY                ClusterID = 0x0405
	CLUSTER_ID_MS_OCCUPANCY_SENSING                ClusterID = 0x0406
	CLUSTER_ID_MS_ELECTRICAL_MEASUREMENT           ClusterID = 0x0b04
)

func (id ClusterID) String() string {
	switch id {
	case CLUSTER_ID_GENERAL_BASIC:
		return "Basic"
	case CLUSTER_ID_GENERAL_POWER_CONFIG:
		return "Power Configuration"
	case CLUSTER_ID_GENERAL_DEVICE_TEMP_CONFIG:
		return "Device Temperature Configuration"
	case CLUSTER_ID_GENERAL_IDENTIFY:
		return "Identify"
	case CLUSTER_ID_GENERAL_GROUPS:
		return "Groups"
	case CLUSTER_ID_GENERAL_SCENES:
		return "Scenes"
	case CLUSTER_ID_GENERAL_ON_OFF:
		return "On/Off"
	case CLUSTER_ID_GENERAL_ON_OFF_SWITCH_CONFIG:
		return "On/Off Switch Configuration"
	case CLUSTER_ID_GENERAL_LEVEL_CONTROL:
		return "Level Control"
	case CLUSTER_ID_GENERAL_ALARMS:
		return "Alarms"
	case CLUSTER_ID_GENERAL_TIME:
		return "Time"
	case CLUSTER_ID_GENERAL_LOCATION:
		return "Location"
	case CLUSTER_ID_GENERAL_DIAGNOSTICS:
		return "Diagnostics"
	case CLUSTER_ID_GENERAL_POLL_CONTROL:
		return "Poll Control"
	case CLUSTER_ID_GENERAL_POWER_PROFILE:
		return "Power Profile"
	case CLUSTER_ID_GENERAL_METER_IDENTIFICATION:
		return "Meter Identification"

	case CLUSTER_ID_GENERAL_ANALOG_INPUT_BASIC:
		return "Analog Input (basic)"
	case CLUSTER_ID_GENERAL_ANALOG_OUTPUT_BASIC:
		return "Analog Output (basic)"
	case CLUSTER_ID_GENERAL_ANALOG_VALUE_BASIC:
		return "Analog Value (basic)"
	case CLUSTER_ID_GENERAL_BINARY_INPUT_BASIC:
		return "Binary Input (basic)"
	case CLUSTER_ID_GENERAL_BINARY_OUTPUT_BASIC:
		return "Binary Output (basic)"
	case CLUSTER_ID_GENERAL_BINARY_VALUE_BASIC:
		return "Binary Value (basic)"
	case CLUSTER_ID_GENERAL_MULTISTATE_INPUT_BASIC:
		return "Multistate Input (basic)"
	case CLUSTER_ID_GENERAL_MULTISTATE_OUTPUT_BASIC:
		return "Multistate Output (basic)"
	case CLUSTER_ID_GENERAL_MULTISTATE_VALUE_BASIC:
		return "Multistate Value (basic)"

	case CLUSTER_ID_MS_ILLUMINANCE_MEASUREMENT:
		return "Illuminance Measurement"
	case CLUSTER_ID_MS_ILLUMINANCE_LEVEL_SENSING_CONFIG:
		return "Illuminance Level Sensing"
	case CLUSTER_ID_MS_TEMPERATURE_MEASUREMENT:
		return "Temperature Measurement"
	case CLUSTER_ID_MS_PRESSURE_MEASUREMENT:
		return "Pressure Measurement"
	case CLUSTER_ID_MS_FLOW_MEASUREMENT:
		return "Flow Measurement"
	case CLUSTER_ID_MS_RELATIVE_HUMIDITY:
		return "Relative Humidity"
	case CLUSTER_ID_MS_OCCUPANCY_SENSING:
		return "Occupancy Sensing"
	case CLUSTER_ID_MS_ELECTRICAL_MEASUREMENT:
		return "Electrical Measurement"

	default:
		return fmt.Sprintf("ClusterID(0x%04x)", uint16(id))
	}
}
