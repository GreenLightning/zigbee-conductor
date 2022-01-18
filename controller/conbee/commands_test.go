package conbee

import (
	"bytes"
	"testing"
)

var cmdtests = []struct {
	name     string
	incoming bool
	data     []byte
}{
	{"ReadFirmwareVersionRequest", false, []byte{0xd, 0x1, 0x0, 0x9, 0x0, 0x0, 0x0, 0xee, 0x2, 0xf9, 0xfe}},
	{"ReadFirmwareVersionResponse", true, []byte{0xd, 0x1, 0x0, 0x9, 0x0, 0x0, 0x7, 0x58, 0x26, 0x64, 0xff}},
	{"ReadParameterRequest", false, []byte{0xa, 0x14, 0x0, 0x8, 0x0, 0x1, 0x0, 0x22, 0xb7, 0xff}},
	{"ReadParameterResponse", true, []byte{0xa, 0x14, 0x0, 0xa, 0x0, 0x3, 0x0, 0x22, 0xb, 0x1, 0xa7, 0xff}},
	{"ReadParameterResponse", true, []byte{0xa, 0x1b, 0x0, 0x10, 0x0, 0x9, 0x0, 0xe, 0x4c, 0x47, 0x7, 0xff, 0xff, 0x2e, 0x21, 0x0, 0xcd, 0xfc}},
	{"WriteParameterRequest", false, []byte{0xb, 0xeb, 0x0, 0x9, 0x0, 0x2, 0x0, 0x21, 0x0, 0xde, 0xfe}},
	{"WriteParameterResponse", true, []byte{0xb, 0xeb, 0x0, 0x8, 0x0, 0x1, 0x0, 0x21, 0xe0, 0xfe}},
	{"DeviceStateRequest", false, []byte{0x7, 0x0, 0x0, 0x8, 0x0, 0x0, 0x0, 0x0, 0xf1, 0xff}},
	{"DeviceStateResponse", true, []byte{0x7, 0x0, 0x0, 0x7, 0x0, 0xa2, 0x0, 0x50, 0xff}},
	{"ReceivedDataNotification", true, []byte{0xe, 0xd7, 0x0, 0x7, 0x0, 0xa6, 0x0, 0x6e, 0xfe}},
	{"ReadReceivedDataRequest", false, []byte{0x17, 0xf2, 0x0, 0x8, 0x0, 0x1, 0x0, 0x6, 0xe8, 0xfe}},
	{"ReadReceivedDataResponse", true, []byte{0x17, 0xf2, 0x0, 0x57, 0x0, 0x50, 0x0, 0x22, 0x2, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x4c, 0x47, 0x7, 0xff, 0xff, 0x2e, 0x21, 0x0, 0x0, 0x0, 0x0, 0x31, 0x80, 0x31, 0x0, 0x7, 0x0, 0x2, 0x0, 0x2, 0x4c, 0x47, 0x7, 0xff, 0xff, 0x2e, 0x21, 0x0, 0x42, 0x50, 0xc5, 0x6, 0x0, 0x8d, 0x15, 0x0, 0x72, 0xe4, 0x12, 0x0, 0x1, 0xff, 0x4c, 0x47, 0x7, 0xff, 0xff, 0x2e, 0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xef, 0x95, 0x25, 0x1, 0x1, 0xff, 0x00, 0x00, 0x11, 0x00, 0x00, 0x00, 0x00, 0x15, 0x4f, 0xed}},
	{"MACPollIndication", true, []byte{0x1c, 0x6d, 0x0, 0x12, 0x0, 0xb, 0x0, 0x3, 0x42, 0x50, 0xc5, 0x6, 0x0, 0x8d, 0x15, 0x0, 0xff, 0xd2, 0x87, 0xfb}},
	{"EnqueueSendDataRequest", false, []byte{0x12, 0xed, 0x0, 0x19, 0x0, 0x12, 0x0, 0xfe, 0x0, 0x2, 0xfc, 0xff, 0x0, 0x0, 0x0, 0x36, 0x0, 0x0, 0x3, 0x0, 0x1e, 0x0, 0x1, 0x0, 0x0, 0x83, 0xfb}},
	{"EnqueueSendDataResponse", true, []byte{0x12, 0xed, 0x0, 0x9, 0x0, 0x2, 0x0, 0x22, 0xfe, 0xd6, 0xfd}},
	{"QuerySendDataRequest", false, []byte{0x4, 0xee, 0x0, 0x7, 0x0, 0x0, 0x0, 0x7, 0xff}},
	{"QuerySendDataResponse", true, []byte{0x4, 0xee, 0x0, 0x13, 0x0, 0xc, 0x0, 0x22, 0xfe, 0x2, 0xfc, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xd2, 0xfb}},
	{"UpdateNeighborCommand", false, []byte{0x1d, 0x3, 0x0, 0x13, 0x0, 0xc, 0x0, 0x1, 0x78, 0x8f, 0x42, 0x50, 0xc5, 0x6, 0x0, 0x8d, 0x15, 0x0, 0x80, 0x3a, 0xfc}},
	{"UpdateNeighborCommand", true, []byte{0x1d, 0x3, 0x0, 0x13, 0x0, 0xc, 0x0, 0x1, 0x78, 0x8f, 0x42, 0x50, 0xc5, 0x6, 0x0, 0x8d, 0x15, 0x0, 0x80, 0x3a, 0xfc}},
}

func TestCommandSerialization(t *testing.T) {
	for _, test := range cmdtests {
		t.Run(test.name, func(t *testing.T) {
			frame, err := ParseFrame(test.data, test.incoming)
			if err != nil {
				t.Fatal("failed to parse frame:", err)
			}
			out, err := SerializeFrame(frame)
			if err != nil {
				t.Fatal("failed to serialize frame:", err)
			}
			if !bytes.Equal(out, test.data) {
				t.Errorf("round-trip error:\ninput:  [% x]\noutput: [% x]\n", test.data, out)
			}
		})
	}
}