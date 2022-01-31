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

	// case ClusterNodeDescReq:
	// case ClusterSimpleDescReq:
	// case ClusterActiveEPReq:
	// case ClusterMgmtLqiReq:

	// case ClusterNWKAddrRsp:
	// case ClusterNodeDescRsp:
	// case ClusterSimpleDescRsp:
	// case ClusterActiveEPRsp:
	// case ClusterMgmtLqiRsp:

	default:
		err = ErrNotImplemented
	}
	return
}

func SerializeFrame(transactionSequenceNumber uint8, command interface{}) (clusterID uint16, data []byte, err error) {
	err = ErrNotImplemented
	return
}
