package zcl

import (
	"reflect"
	"testing"
)

func TestParseValueSuccess(t *testing.T) {
	type TestCase struct {
		DataType DataType
		Input    []byte
		Output   interface{}
	}

	testCases := []TestCase{
		TestCase{DataTypeNoData, []byte{42}, nil},

		TestCase{DataTypeData24, []byte{1, 2, 3, 42}, []byte{1, 2, 3}},

		TestCase{DataTypeBool, []byte{0x00, 42}, false},
		TestCase{DataTypeBool, []byte{0x01, 42}, true},
		TestCase{DataTypeBool, []byte{0xff, 42}, nil},

		TestCase{DataTypeBitmap8, []byte{0xaa, 42}, uint8(0xaa)},
		TestCase{DataTypeBitmap16, []byte{0xbb, 0xaa, 42}, uint16(0xaabb)},
		TestCase{DataTypeBitmap24, []byte{0xcc, 0xbb, 0xaa, 42}, uint32(0xaabbcc)},
		TestCase{DataTypeBitmap32, []byte{0xdd, 0xcc, 0xbb, 0xaa, 42}, uint32(0xaabbccdd)},
		TestCase{DataTypeBitmap40, []byte{0xee, 0xdd, 0xcc, 0xbb, 0xaa, 42}, uint64(0xaabbccddee)},
		TestCase{DataTypeBitmap48, []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 42}, uint64(0xaabbccddeeff)},
		TestCase{DataTypeBitmap56, []byte{0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 42}, uint64(0xaabbccddeeff00)},
		TestCase{DataTypeBitmap64, []byte{0x11, 0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 42}, uint64(0xaabbccddeeff0011)},

		TestCase{DataTypeBitmap24, []byte{0xff, 0xff, 0xff, 42}, uint32(0xffffff)},

		TestCase{DataTypeUint8, []byte{0xaa, 42}, uint8(0xaa)},
		TestCase{DataTypeUint16, []byte{0xbb, 0xaa, 42}, uint16(0xaabb)},
		TestCase{DataTypeUint24, []byte{0xcc, 0xbb, 0xaa, 42}, uint32(0xaabbcc)},
		TestCase{DataTypeUint32, []byte{0xdd, 0xcc, 0xbb, 0xaa, 42}, uint32(0xaabbccdd)},
		TestCase{DataTypeUint40, []byte{0xee, 0xdd, 0xcc, 0xbb, 0xaa, 42}, uint64(0xaabbccddee)},
		TestCase{DataTypeUint48, []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 42}, uint64(0xaabbccddeeff)},
		TestCase{DataTypeUint56, []byte{0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 42}, uint64(0xaabbccddeeff00)},
		TestCase{DataTypeUint64, []byte{0x11, 0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 42}, uint64(0xaabbccddeeff0011)},

		TestCase{DataTypeUint8, []byte{0xff, 42}, nil},
		TestCase{DataTypeUint16, []byte{0xff, 0xff, 42}, nil},
		TestCase{DataTypeUint24, []byte{0xff, 0xff, 0xff, 42}, nil},
		TestCase{DataTypeUint32, []byte{0xff, 0xff, 0xff, 0xff, 42}, nil},
		TestCase{DataTypeUint40, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 42}, nil},
		TestCase{DataTypeUint48, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 42}, nil},
		TestCase{DataTypeUint56, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 42}, nil},
		TestCase{DataTypeUint64, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 42}, nil},

		TestCase{DataTypeInt8, []byte{0x01, 42}, int8(0x01)},
		TestCase{DataTypeInt16, []byte{0x02, 0x01, 42}, int16(0x0102)},
		TestCase{DataTypeInt24, []byte{0x03, 0x02, 0x01, 42}, int32(0x010203)},
		TestCase{DataTypeInt32, []byte{0x04, 0x03, 0x02, 0x01, 42}, int32(0x01020304)},
		TestCase{DataTypeInt40, []byte{0x05, 0x04, 0x03, 0x02, 0x01, 42}, int64(0x0102030405)},
		TestCase{DataTypeInt48, []byte{0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 42}, int64(0x010203040506)},
		TestCase{DataTypeInt56, []byte{0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 42}, int64(0x01020304050607)},
		TestCase{DataTypeInt64, []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 42}, int64(0x0102030405060708)},

		TestCase{DataTypeInt8, []byte{0xff, 42}, int8(-1)},
		TestCase{DataTypeInt16, []byte{0xff, 0xff, 42}, int16(-1)},
		TestCase{DataTypeInt24, []byte{0xff, 0xff, 0xff, 42}, int32(-1)},
		TestCase{DataTypeInt32, []byte{0xff, 0xff, 0xff, 0xff, 42}, int32(-1)},
		TestCase{DataTypeInt40, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 42}, int64(-1)},
		TestCase{DataTypeInt48, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 42}, int64(-1)},
		TestCase{DataTypeInt56, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 42}, int64(-1)},
		TestCase{DataTypeInt64, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 42}, int64(-1)},

		TestCase{DataTypeInt8, []byte{0x80, 42}, nil},
		TestCase{DataTypeInt16, []byte{0x00, 0x80, 42}, nil},
		TestCase{DataTypeInt24, []byte{0x00, 0x00, 0x80, 42}, nil},
		TestCase{DataTypeInt32, []byte{0x00, 0x00, 0x00, 0x80, 42}, nil},
		TestCase{DataTypeInt40, []byte{0x00, 0x00, 0x00, 0x00, 0x80, 42}, nil},
		TestCase{DataTypeInt48, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 42}, nil},
		TestCase{DataTypeInt56, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 42}, nil},
		TestCase{DataTypeInt64, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 42}, nil},

		TestCase{DataTypeFloat32, []byte{0x00, 0x00, 0x00, 0x3e, 42}, float32(0.125)},
		TestCase{DataTypeFloat64, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x3f, 42}, float64(0.125)},

		TestCase{DataTypeFloat32, []byte{0x01, 0x00, 0x80, 0x7f, 42}, nil},
		TestCase{DataTypeFloat64, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x7f, 42}, nil},

		TestCase{DataTypeCharacterString, []byte{0x05, 'H', 'e', 'l', 'l', 'o', 42}, "Hello"},
		TestCase{DataTypeCharacterString, []byte{0xff, 42}, nil},

		TestCase{DataTypeLongCharacterString, []byte{0x05, 0x00, 'H', 'e', 'l', 'l', 'o', 42}, "Hello"},
		TestCase{DataTypeLongCharacterString, []byte{0xff, 0xff, 42}, nil},
	}

	for index, testCase := range testCases {
		value, data, err := ParseValue(testCase.DataType, testCase.Input)

		if err != nil {
			t.Errorf("(%d) unexpected err: %v", index, err)
			continue
		}

		if len(data) != 1 || data[0] != 42 {
			t.Errorf("(%d) wrong data: %v", index, data)
		}

		if reflect.TypeOf(value) != reflect.TypeOf(testCase.Output) {
			t.Errorf("(%d) wrong value type: expected %v (%v): %v (%v)", index, reflect.TypeOf(testCase.Output), testCase.Output, reflect.TypeOf(value), value)
			continue
		}

		// Some special handling, because we cannot compare slices using "=="!
		if typ := reflect.TypeOf(testCase.Output); typ != nil && typ.Kind() == reflect.Slice {
			output2 := reflect.ValueOf(testCase.Output)
			value2 := reflect.ValueOf(value)
			if value2.Len() != output2.Len() {
				t.Errorf("(%d) wrong value (slice length): expected %v: %v", index, testCase.Output, value)
				continue
			}

			for i := 0; i < output2.Len(); i++ {
				if value2.Index(i).Interface() != output2.Index(i).Interface() {
					t.Errorf("(%d) wrong value (slice contents): expected %v: %v", index, testCase.Output, value)
					break
				}
			}
		} else {
			if value != testCase.Output {
				t.Errorf("(%d) wrong value (simple): expected %v: %v", index, testCase.Output, value)
			}
		}
	}
}

func TestParseValueDataNotEnoughData(t *testing.T) {
	_, _, err := ParseValue(DataTypeData24, []byte{1, 2})

	if err != ErrNotEnoughData {
		t.Fatal("expected ErrNotEnoughData:", err)
	}
}

func TestParseValueBoolInvalidData(t *testing.T) {
	_, _, err := ParseValue(DataTypeBool, []byte{2})

	if err != ErrInvalidData {
		t.Fatal("expected ErrInvalidData:", err)
	}
}

func TestParseValueCharacterStringNotEnoughData1(t *testing.T) {
	_, _, err := ParseValue(DataTypeCharacterString, []byte{})

	if err != ErrNotEnoughData {
		t.Fatal("expected ErrNotEnoughData:", err)
	}
}

func TestParseValueCharacterStringNotEnoughData2(t *testing.T) {
	_, _, err := ParseValue(DataTypeCharacterString, []byte{1})

	if err != ErrNotEnoughData {
		t.Fatal("expected ErrNotEnoughData:", err)
	}
}

func TestParseValueLongCharacterStringNotEnoughData1(t *testing.T) {
	_, _, err := ParseValue(DataTypeLongCharacterString, []byte{1})

	if err != ErrNotEnoughData {
		t.Fatal("expected ErrNotEnoughData:", err)
	}
}

func TestParseValueLongCharacterStringNotEnoughData2(t *testing.T) {
	_, _, err := ParseValue(DataTypeLongCharacterString, []byte{1, 0})

	if err != ErrNotEnoughData {
		t.Fatal("expected ErrNotEnoughData:", err)
	}
}
