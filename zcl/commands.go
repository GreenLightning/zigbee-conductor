package zcl

import (
	"encoding/binary"
	"fmt"
)

const (
	COMMAND_ID_READ_ATTRIBUTES                       = 0x00
	COMMAND_ID_READ_ATTRIBUTES_RESPONSE              = 0x01
	COMMAND_ID_WRITE_ATTRIBUTES                      = 0x02
	COMMAND_ID_WRITE_ATTRIBUTES_UNDIVIDED            = 0x03
	COMMAND_ID_WRITE_ATTRIBUTES_RESPONSE             = 0x04
	COMMAND_ID_WRITE_ATTRIBUTES_NO_RESPONSE          = 0x05
	COMMAND_ID_CONFIGURE_REPORTING                   = 0x06
	COMMAND_ID_CONFIGURE_REPORTING_RESPONSE          = 0x07
	COMMAND_ID_READ_REPORTING_CONFIGURATION          = 0x08
	COMMAND_ID_READ_REPORTING_CONFIGURATION_RESPONSE = 0x09
	COMMAND_ID_REPORT_ATTRIBUTES                     = 0x0a
	COMMAND_ID_DEFAULT_RESPONSE                      = 0x0b
	COMMAND_ID_DISCOVER_ATTRIBUTES                   = 0x0c
	COMMAND_ID_DISCOVER_ATTRIBUTES_RESPONSE          = 0x0d
	COMMAND_ID_READ_ATTRIBUTES_STRUCTURED            = 0x0e
	COMMAND_ID_WRITE_ATTRIBUTES_STRUCTURED           = 0x0f
	COMMAND_ID_WRITE_ATTRIBUTES_STRUCTURED_RESPONSE  = 0x10
	COMMAND_ID_DISCOVER_COMMANDS_RECEIVED            = 0x11
	COMMAND_ID_DISCOVER_COMMANDS_RECEIVED_RESPONSE   = 0x12
	COMMAND_ID_DISCOVER_COMMANDS_GENERATED           = 0x13
	COMMAND_ID_DISCOVER_COMMANDS_GENERATED_RESPONSE  = 0x14
	COMMAND_ID_DISCOVER_ATTRIBUTES_EXTENDED          = 0x15
	COMMAND_ID_DISCOVER_ATTRIBUTES_EXTENDED_RESPONSE = 0x16
)

type Status byte

const (
	StatusSuccess                               Status = 0x00
	StatusFailure                               Status = 0x01
	StatusNotAuthorized                         Status = 0x7e
	StatusReservedFieldNotZero                  Status = 0x7f
	StatusMalformedCommand                      Status = 0x80
	StatusUnsupportedClusterCommand             Status = 0x81
	StatusUnsupportedGeneralCommand             Status = 0x82
	StatusUnsupportedManufacturerClusterCommand Status = 0x83
	StatusUnsupportedManufacturerGeneralCommand Status = 0x84
	StatusInvalidField                          Status = 0x85
	StatusUnsupportedAttribute                  Status = 0x86
	StatusInvalidValue                          Status = 0x87
	StatusReadOnly                              Status = 0x88
	StatusInsufficientSpace                     Status = 0x89
	StatusDuplicateExists                       Status = 0x8a
	StatusNotFound                              Status = 0x8b
	StatusUnreportableAttribute                 Status = 0x8c
	StatusInvalidDataType                       Status = 0x8d
	StatusInvalidSelector                       Status = 0x8e
	StatusWriteOnly                             Status = 0x8f
	StatusInconsistentStartupState              Status = 0x90
	StatusDefinedOutOfBand                      Status = 0x91
	StatusInconsistent                          Status = 0x92
	StatusActionDenied                          Status = 0x93
	StatusTimeout                               Status = 0x94
	StatusAbort                                 Status = 0x95
	StatusInvalidImage                          Status = 0x96
	StatusWaitForData                           Status = 0x97
	StatusNoImageAvailable                      Status = 0x98
	StatusRequireMoreImage                      Status = 0x99
	StatusNotificationPending                   Status = 0x9a
	StatusHardwareFailure                       Status = 0xc0
	StatusSoftwareFailure                       Status = 0xc1
	StatusCalibrationError                      Status = 0xc2
	StatusUnsupportedCluster                    Status = 0xc3
)

func (status Status) String() string {
	switch status {
	case StatusSuccess:
		return "Success"
	case StatusFailure:
		return "Failure"
	case StatusNotAuthorized:
		return "NotAuthorized"
	case StatusReservedFieldNotZero:
		return "ReservedFieldNotZero"
	case StatusMalformedCommand:
		return "MalformedCommand"
	case StatusUnsupportedClusterCommand:
		return "UnsupportedClusterCommand"
	case StatusUnsupportedGeneralCommand:
		return "UnsupportedGeneralCommand"
	case StatusUnsupportedManufacturerClusterCommand:
		return "UnsupportedManufacturerClusterCommand"
	case StatusUnsupportedManufacturerGeneralCommand:
		return "UnsupportedManufacturerGeneralCommand"
	case StatusInvalidField:
		return "InvalidField"
	case StatusUnsupportedAttribute:
		return "UnsupportedAttribute"
	case StatusInvalidValue:
		return "InvalidValue"
	case StatusReadOnly:
		return "ReadOnly"
	case StatusInsufficientSpace:
		return "InsufficientSpace"
	case StatusDuplicateExists:
		return "DuplicateExists"
	case StatusNotFound:
		return "NotFound"
	case StatusUnreportableAttribute:
		return "UnreportableAttribute"
	case StatusInvalidDataType:
		return "InvalidDataType"
	case StatusInvalidSelector:
		return "InvalidSelector"
	case StatusWriteOnly:
		return "WriteOnly"
	case StatusInconsistentStartupState:
		return "InconsistentStartupState"
	case StatusDefinedOutOfBand:
		return "DefinedOutOfBand"
	case StatusInconsistent:
		return "Inconsistent"
	case StatusActionDenied:
		return "ActionDenied"
	case StatusTimeout:
		return "Timeout"
	case StatusAbort:
		return "Abort"
	case StatusInvalidImage:
		return "InvalidImage"
	case StatusWaitForData:
		return "WaitForData"
	case StatusNoImageAvailable:
		return "NoImageAvailable"
	case StatusRequireMoreImage:
		return "RequireMoreImage"
	case StatusNotificationPending:
		return "NotificationPending"
	case StatusHardwareFailure:
		return "HardwareFailure"
	case StatusSoftwareFailure:
		return "SoftwareFailure"
	case StatusCalibrationError:
		return "CalibrationError"
	case StatusUnsupportedCluster:
		return "UnsupportedCluster"
	default:
		return fmt.Sprintf("Status(0x%02x)", byte(status))
	}
}

type AttributeID uint16

func (id AttributeID) String() string {
	return fmt.Sprintf("0x%04x", uint16(id))
}

type ReadReportingConfigurationCommand struct {
	Records []AttributeRecord
}

type AttributeRecord struct {
	Direction   byte
	AttributeID AttributeID
}

func SerializeReadReportingConfigurationCommand(command ReadReportingConfigurationCommand) []byte {
	data := make([]byte, 0, 3*len(command.Records))

	for _, record := range command.Records {
		data = append(data, record.Direction, 0, 0)
		binary.LittleEndian.PutUint16(data[len(data)-2:], uint16(record.AttributeID))
	}

	return data
}

type ReadReportingConfigurationResponseCommand struct {
	Records []AttributeReportingConfigurationRecord
}

type AttributeReportingConfigurationRecord struct {
	Status                   Status
	Direction                byte
	AttributeID              AttributeID
	DataType                 DataType
	MinimumReportingInterval uint16
	MaximumReportingInterval uint16
	ReportableChange         interface{}
	TimeoutPeriod            uint16
}

func ParseReadReportingConfigurationResponseCommand(data []byte) (ReadReportingConfigurationResponseCommand, error) {
	var command ReadReportingConfigurationResponseCommand
	for len(data) != 0 {
		var record AttributeReportingConfigurationRecord
		var err error
		record, data, err = ParseAttributeReportingConfigurationRecord(data)
		if err != nil {
			return command, err
		}
		command.Records = append(command.Records, record)
	}
	return command, nil
}

func ParseAttributeReportingConfigurationRecord(data []byte) (AttributeReportingConfigurationRecord, []byte, error) {
	var record AttributeReportingConfigurationRecord

	if len(data) < 4 {
		return record, data, ErrNotEnoughData
	}

	record.Status = Status(data[0])
	record.Direction = data[1]
	data = data[2:]

	record.AttributeID = AttributeID(binary.LittleEndian.Uint16(data))
	data = data[2:]

	if record.Status == StatusSuccess {
		if len(data) < 5 {
			return record, data, ErrNotEnoughData
		}

		record.DataType = DataType(data[0])
		data = data[1:]

		record.MinimumReportingInterval = binary.LittleEndian.Uint16(data)
		data = data[2:]

		record.MaximumReportingInterval = binary.LittleEndian.Uint16(data)
		data = data[2:]

		var err error
		record.ReportableChange, data, err = ParseValue(record.DataType, data)
		if err != nil {
			return record, data, err
		}

		if len(data) < 2 {
			return record, data, ErrNotEnoughData
		}

		record.TimeoutPeriod = binary.LittleEndian.Uint16(data)
		data = data[2:]
	}

	return record, data, nil
}

type ReportAttributesCommand struct {
	Reports []AttributeReport
}

type AttributeReport struct {
	AttributeID AttributeID
	DataType    DataType
	Value       interface{}
}

func ParseReportAttributesCommand(data []byte) (ReportAttributesCommand, error) {
	var command ReportAttributesCommand
	for len(data) != 0 {
		var report AttributeReport
		var err error
		report, data, err = ParseAttributeReport(data)
		if err != nil {
			return command, err
		}
		command.Reports = append(command.Reports, report)
	}
	return command, nil
}

func ParseAttributeReport(data []byte) (AttributeReport, []byte, error) {
	var report AttributeReport

	if len(data) < 3 {
		return report, data, ErrNotEnoughData
	}

	report.AttributeID = AttributeID(binary.LittleEndian.Uint16(data))
	data = data[2:]

	report.DataType = DataType(data[0])
	data = data[1:]

	var err error
	report.Value, data, err = ParseValue(report.DataType, data)
	return report, data, err
}
