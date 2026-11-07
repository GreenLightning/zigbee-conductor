package zcl

import (
	"fmt"
	"strings"
	"testing"
)

func TestBinaryReportFormatting(t *testing.T) {
	r := AttributeReport{AttributeID: 0xff02, DataType: DataTypeStructure, Value: []TypedValue{{DataType: DataTypeCharacterString, Value: "\x00\n\xff"}}}
	got := fmt.Sprintf("%+v", r)
	if strings.ContainsAny(got, "\x00\n\r") || !strings.Contains(got, "hex:000aff") {
		t.Fatal(got)
	}
}
