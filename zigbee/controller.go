package zigbee

import "io"

type ControllerSettings struct {
	Port        string
	LogCommands bool
	LogErrors   bool
}

type Controller interface {
	io.Closer

	Start() (chan IncomingMessage, error)
	Send(msg OutgoingMessage) error
	PermitJoining(enabled bool) error
}

type IncomingMessage struct {
	Source              Address
	SourceEndpoint      uint8
	DestinationEndpoint uint8
	ClusterID           uint16
	LinkQuality         uint8
	Data                []byte
}

type OutgoingMessage struct {
	Destination         Address
	DestinationEndpoint uint8
	SourceEndpoint      uint8
	ClusterID           uint16
	Radius              uint8
	Data                []byte
}
