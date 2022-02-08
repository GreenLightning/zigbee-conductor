package zdp

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

// Note: For testing the transaction sequence number is always 0xAB.

type TestCase struct {
	ClusterID uint16
	Data      string
	Command   interface{}
}

var tests = []TestCase{
	TestCase{ClusterID: 0x0000, Data: "AB77665544330200AA0001", Command: &NWKAddrReq{IEEEAddress: 0xAA00023344556677, RequestType: 0, StartIndex: 1}},
}

func TestParseFrame(t *testing.T) {
	for _, tc := range tests {
		t.Run(ClusterName(tc.ClusterID), func(t *testing.T) {
			data, _ := hex.DecodeString(tc.Data)
			tsn, cmd, err := ParseFrame(tc.ClusterID, data)
			if err != nil {
				t.Fatal("unexpected err:", err)
			}
			if tsn != 0xAB {
				t.Errorf("wrong TSN: expected 0x%02x, actual 0x%02x", 0xAB, tsn)
			}
			if !reflect.DeepEqual(cmd, tc.Command) {
				t.Errorf("wrong command:\n\texpected %T%+v\n\tactual   %T%+v", tc.Command, tc.Command, cmd, cmd)
			}
		})
	}
}

func TestSerializeFrame(t *testing.T) {
	for _, tc := range tests {
		t.Run(ClusterName(tc.ClusterID), func(t *testing.T) {
			expected, _ := hex.DecodeString(tc.Data)
			clusterID, actual, err := SerializeFrame(0xAB, tc.Command)
			if err != nil {
				t.Fatal("unexpected err:", err)
			}
			if clusterID != tc.ClusterID {
				t.Errorf("wrong cluster ID: expected 0x%04x, actual 0x%04x", tc.ClusterID, clusterID)
			}
			if !bytes.Equal(actual, expected) {
				t.Errorf("wrong data:\n\texpected %x\n\tactual   %x", expected, actual)
			}
		})
	}
}
