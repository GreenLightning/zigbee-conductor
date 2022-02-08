// Package zdp implements the ZigBee Device Profile.
package zdp

import (
	"errors"

	"github.com/GreenLightning/zigbee-conductor/pkg/scf"
)

var (
	ErrInvalidData    = errors.New("invalid data")
	ErrNotImplemented = errors.New("not implemented")
)

func ParseFrame(clusterID uint16, data []byte) (transactionSequenceNumber uint8, command interface{}, err error) {
	if len(data) == 0 {
		err = ErrInvalidData
		return
	}
	transactionSequenceNumber = data[0]
	data = data[1:]

	switch clusterID {
	case ClusterNWKAddrReq:
		command = new(NWKAddrReq)
		_, err = scf.Parse(command, data)

	// case ClusterIEEEAddrReq:
	// case ClusterActiveEPReq:
	// case ClusterMatchDescReq:

	case ClusterDeviceAnnce:
		command = new(DeviceAnnce)
		_, err = scf.Parse(command, data)

	// case ClusterParentAnnce:
	// case ClusterEndDeviceBindReq:
	// case ClusterBindReq:
	// case ClusterUnbindReq:
	// case ClusterBindRegisterReq:
	// case ClusterReplaceDeviceReq:

	// case ClusterNWKAddrRsp:
	// case ClusterIEEEAddrRsp:
	// case ClusterActiveEPRsp:
	// case ClusterMatchDescRsp:
	// case ClusterParentAnnceRsp:
	// case ClusterEndDeviceBindRsp:
	// case ClusterBindRsp:
	// case ClusterUnbindRsp:
	// case ClusterBindRegisterRsp:
	// case ClusterReplaceDeviceRsp:

	default:
		err = ErrNotImplemented
	}
	return
}

func SerializeFrame(transactionSequenceNumber uint8, command interface{}) (clusterID uint16, data []byte, err error) {
	data = append(data, transactionSequenceNumber)
	switch  command := command.(type) {
	case *NWKAddrReq:
		clusterID = ClusterNWKAddrReq
		data = append(data, scf.Serialize(*command)...)

	// case *IEEEAddrReq:
	// case *ActiveEPReq:
	// case *MatchDescReq:
	// case *DeviceAnnce:
	// case *ParentAnnce:
	// case *EndDeviceBindReq:
	// case *BindReq:
	// case *UnbindReq:
	// case *BindRegisterReq:
	// case *ReplaceDeviceReq:
	// case *NWKAddrRsp:
	// case *IEEEAddrRsp:
	// case *ActiveEPRsp:
	// case *MatchDescRsp:
	// case *ParentAnnceRsp:
	// case *EndDeviceBindRsp:
	// case *BindRsp:
	// case *UnbindRsp:
	// case *BindRegisterRsp:
	// case *ReplaceDeviceRsp:

	default:
		err = ErrNotImplemented
	}
	return
}
