package zcl

import "fmt"

type ClusterID uint16

const (
	ClusterGeneralBasic               ClusterID = 0x0000
	ClusterGeneralPowerConfig         ClusterID = 0x0001
	ClusterGeneralDeviceTempConfig    ClusterID = 0x0002
	ClusterGeneralIdentify            ClusterID = 0x0003
	ClusterGeneralGroups              ClusterID = 0x0004
	ClusterGeneralScenes              ClusterID = 0x0005
	ClusterGeneralOnOff               ClusterID = 0x0006
	ClusterGeneralOnOffSwitchConfig   ClusterID = 0x0007
	ClusterGeneralLevelControl        ClusterID = 0x0008
	ClusterGeneralAlarms              ClusterID = 0x0009
	ClusterGeneralTime                ClusterID = 0x000a
	ClusterGeneralLocation            ClusterID = 0x000b
	ClusterGeneralDiagnostics         ClusterID = 0x0b05
	ClusterGeneralPollControl         ClusterID = 0x0020
	ClusterGeneralPowerProfile        ClusterID = 0x001a
	ClusterGeneralMeterIdentification ClusterID = 0x0b01

	ClusterGeneralAnalogInputBasic      ClusterID = 0x000c
	ClusterGeneralAnalogOutputBasic     ClusterID = 0x000d
	ClusterGeneralAnalogValueBasic      ClusterID = 0x000e
	ClusterGeneralBinaryInputBasic      ClusterID = 0x000f
	ClusterGeneralBinaryOutputBasic     ClusterID = 0x0010
	ClusterGeneralBinaryValueBasic      ClusterID = 0x0011
	ClusterGeneralMultistateInputBasic  ClusterID = 0x0012
	ClusterGeneralMultistateOutputBasic ClusterID = 0x0013
	ClusterGeneralMultistateValueBasic  ClusterID = 0x0014

	ClusterMSIlluminanceMeasurement        ClusterID = 0x0400
	ClusterMSIlluminanceLevelSensingConfig ClusterID = 0x0401
	ClusterMSTemperatureMeasurement        ClusterID = 0x0402
	ClusterMSPressureMeasurement           ClusterID = 0x0403
	ClusterMSFlowMeasurement               ClusterID = 0x0404
	ClusterMSRelativeHumidity              ClusterID = 0x0405
	ClusterMSOccupancySensing              ClusterID = 0x0406
	ClusterMSElectricalMeasurement         ClusterID = 0x0b04
)

func (id ClusterID) String() string {
	switch id {
	case ClusterGeneralBasic:
		return "Basic"
	case ClusterGeneralPowerConfig:
		return "PowerConfig"
	case ClusterGeneralDeviceTempConfig:
		return "DeviceTempConfig"
	case ClusterGeneralIdentify:
		return "Identify"
	case ClusterGeneralGroups:
		return "Groups"
	case ClusterGeneralScenes:
		return "Scenes"
	case ClusterGeneralOnOff:
		return "OnOff"
	case ClusterGeneralOnOffSwitchConfig:
		return "OnOffSwitchConfig"
	case ClusterGeneralLevelControl:
		return "LevelControl"
	case ClusterGeneralAlarms:
		return "Alarms"
	case ClusterGeneralTime:
		return "Time"
	case ClusterGeneralLocation:
		return "Location"
	case ClusterGeneralDiagnostics:
		return "Diagnostics"
	case ClusterGeneralPollControl:
		return "PollControl"
	case ClusterGeneralPowerProfile:
		return "PowerProfile"
	case ClusterGeneralMeterIdentification:
		return "MeterIdentification"

	case ClusterGeneralAnalogInputBasic:
		return "AnalogInputBasic"
	case ClusterGeneralAnalogOutputBasic:
		return "AnalogOutputBasic"
	case ClusterGeneralAnalogValueBasic:
		return "AnalogValueBasic"
	case ClusterGeneralBinaryInputBasic:
		return "BinaryInputBasic"
	case ClusterGeneralBinaryOutputBasic:
		return "BinaryOutputBasic"
	case ClusterGeneralBinaryValueBasic:
		return "BinaryValueBasic"
	case ClusterGeneralMultistateInputBasic:
		return "MultistateInputBasic"
	case ClusterGeneralMultistateOutputBasic:
		return "MultistateOutputBasic"
	case ClusterGeneralMultistateValueBasic:
		return "MultistateValueBasic"

	case ClusterMSIlluminanceMeasurement:
		return "IlluminanceMeasurement"
	case ClusterMSIlluminanceLevelSensingConfig:
		return "IlluminanceLevelSensingConfig"
	case ClusterMSTemperatureMeasurement:
		return "TemperatureMeasurement"
	case ClusterMSPressureMeasurement:
		return "PressureMeasurement"
	case ClusterMSFlowMeasurement:
		return "FlowMeasurement"
	case ClusterMSRelativeHumidity:
		return "RelativeHumidity"
	case ClusterMSOccupancySensing:
		return "OccupancySensing"
	case ClusterMSElectricalMeasurement:
		return "ElectricalMeasurement"

	default:
		return fmt.Sprintf("ClusterID(0x%04x)", uint16(id))
	}
}
