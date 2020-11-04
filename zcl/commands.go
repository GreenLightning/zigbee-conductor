package zcl

import (
	"encoding/binary"
	"errors"
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

type AttributeID uint16

func (id AttributeID) String() string {
	return fmt.Sprintf("0x%04x", uint16(id))
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
		return report, data, errors.New("not enough data")
	}

	report.AttributeID = AttributeID(binary.LittleEndian.Uint16(data))
	data = data[2:]

	report.DataType = DataType(data[0])
	data = data[1:]

	var err error
	report.Value, data, err = ParseValue(report.DataType, data)
	return report, data, err
}
