package zcl

import (
	"encoding/hex"
	"testing"
)

func TestParseCapturedStructure(t *testing.T) {
	data, _ := hex.DecodeString("0600100121c40b21a843240100000000211503205d")
	value, rest, err := ParseValue(DataTypeStructure, append(data, 0xaa))
	if err != nil {
		t.Fatal(err)
	}
	elements := value.([]TypedValue)
	if len(elements) != 6 || elements[1].Value != uint16(3012) || len(rest) != 1 || rest[0] != 0xaa {
		t.Fatalf("unexpected structure: %+v rest=%x", elements, rest)
	}
	for n := 0; n < len(data); n++ {
		if _, _, err := ParseValue(DataTypeStructure, data[:n]); err == nil {
			t.Errorf("accepted truncation at %d", n)
		}
	}
}
func TestParseNestedStructure(t *testing.T) {
	value, _, err := ParseValue(DataTypeStructure, []byte{1, 0, 0x4c, 1, 0, 0x21, 0x34, 0x12})
	if err != nil {
		t.Fatal(err)
	}
	if value.([]TypedValue)[0].Value.([]TypedValue)[0].Value != uint16(0x1234) {
		t.Fatal(value)
	}
}
