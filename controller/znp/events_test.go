package znp

import (
	"github.com/GreenLightning/zigbee-conductor/zigbee"
	"testing"
)

func TestDeviceEvents(t *testing.T) {
	for _, tt := range []struct {
		command interface{}
		kind    zigbee.DeviceEventKind
	}{
		{ZdoEndDeviceAnnceInd{NwkAddr: 42, IEEEAddr: 1234}, zigbee.DeviceAnnounced},
		{ZdoTcDevInd{SrcNwkAddr: 42, SrcIEEEAddr: 1234, ParentNwkAddr: 7}, zigbee.DeviceJoined},
	} {
		e, ok := deviceEvent(tt.command)
		if !ok || e.Kind != tt.kind || e.NetworkAddress != 42 || e.IEEEAddress != 1234 {
			t.Fatalf("unexpected event: %+v", e)
		}
	}
	if _, ok := deviceEvent(SysVersionRequest{}); ok {
		t.Fatal("non-event accepted")
	}
}
