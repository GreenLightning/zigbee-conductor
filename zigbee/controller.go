package zigbee

import "io"

type ControllerSettings struct {
	Port        string
	LogCommands bool
	LogErrors   bool
}

// A controller allows interacting with the ZigBee network on the application level.
//
// The most important feature is sending and receiving messages (more
// specifically these messages are based on the Application Support (APS) layer,
// containing the APS payload and fields from the APS layer (e.g. cluster and
// profile ID) and lower layers (e.g. addresses from the Network layer)).
//
// Endpoint management is still an open issue. You can expect endpoint 1 to be
// present and configured using the Home Automation profile. If you require more
// advanced endpoint management, please open an issue.
type Controller interface {
	io.Closer

	Start() (chan IncomingMessage, error)
	Send(message OutgoingMessage) error
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
