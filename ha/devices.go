// Package ha contains constants from the Home Automation profile.
package ha

// Device IDs specified in the Home Automation profile.
// Extracted from Home Automation Public Application Profile Specification Revision 29, Version 1.2 (2013-06-06).

// Generic
const (
	DeviceOnOffSwitch                = 0x0000
	DeviceLevelControlSwitch         = 0x0001
	DeviceOnOffOutput                = 0x0002
	DeviceLevelControllableOutput    = 0x0003
	DeviceSceneSelector              = 0x0004
	DeviceConfigurationTool          = 0x0005
	DeviceRemoteControl              = 0x0006
	DeviceCombinedInterface          = 0x0007
	DeviceRangeExtender              = 0x0008
	DeviceMainsPowerOutlet           = 0x0009
	DeviceDoorLock                   = 0x000A
	DeviceDoorLockController         = 0x000B
	DeviceSimpleSensor               = 0x000C
	DeviceConsumptionAwarenessDevice = 0x000D
	DeviceHomeGateway                = 0x0050
	DeviceSmartPlug                  = 0x0051
	DeviceWhiteGoods                 = 0x0052
	DeviceMeterInterface             = 0x0053
)

// Lighting
const (
	DeviceOnOffLight         = 0x0100
	DeviceDimmableLight      = 0x0101
	DeviceColorDimmableLight = 0x0102
	DeviceOnOffLightSwitch   = 0x0103
	DeviceDimmerSwitch       = 0x0104
	DeviceColorDimmerSwitch  = 0x0105
	DeviceLightSensor        = 0x0106
	DeviceOccupancySensor    = 0x0107
)

// Closures
const (
	DeviceClosuresShade            = 0x0200
	DeviceShadeController          = 0x0201
	DeviceWindowCoveringDevice     = 0x0202
	DeviceWindowCoveringController = 0x0203
)

// HVAC
const (
	DeviceHeatingCoolingUnit = 0x0300
	DeviceThermostat         = 0x0301
	DeviceTemperatureSensor  = 0x0302
	DevicePump               = 0x0303
	DevicePumpController     = 0x0304
	DevicePressureSensor     = 0x0305
	DeviceFlowSensor         = 0x0306
	DeviceMiniSplitAC        = 0x0307
)

// Intruder Alarm Systems
const (
	DeviceIASControlAndIndicatingEquipment = 0x0400
	DeviceIASAncillaryControlEquipment     = 0x0401
	DeviceIASZone                          = 0x0402
	DeviceIASWarningDevice                 = 0x0403
)
