package zcl

import (
	"fmt"
	"strings"
)

// FormatValue renders binary strings and byte arrays as hex, including nested structures.
func FormatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("hex:%x", []byte(v))
	case []byte:
		return fmt.Sprintf("hex:%x", v)
	case []TypedValue:
		parts := make([]string, 0, len(v))
		for _, element := range v {
			parts = append(parts, fmt.Sprintf("{DataType:%v Value:%s}", element.DataType, FormatValue(element.Value)))
		}
		return "[" + strings.Join(parts, " ") + "]"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (report AttributeReport) String() string {
	return fmt.Sprintf("{AttributeID:0x%04x DataType:%v Value:%s}", uint16(report.AttributeID), report.DataType, FormatValue(report.Value))
}
