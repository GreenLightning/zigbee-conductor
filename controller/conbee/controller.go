// Implements a Controller for the Phoscon ConBee II.
//
// The serial protocol is documented here:
//
// https://github.com/dresden-elektronik/deconz-serial-protocol
//
// https://deconz.dresden-elektronik.de/raspbian/deCONZ-Serial-Protocol-en_1.21.pdf
package conbee

import (
	"fmt"
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/GreenLightning/zigbee-conductor/pkg/slip"
	"github.com/GreenLightning/zigbee-conductor/zigbee"
	"github.com/jacobsa/go-serial/serial"
)

type Controller struct {
	settings        zigbee.ControllerSettings
	sequence        uint32
	requestSequence uint32
	port            io.ReadWriteCloser
}

func NewController(settings zigbee.ControllerSettings) (*Controller, error) {
	options := serial.OpenOptions{
		PortName:        settings.Port,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	port, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	controller := &Controller{
		settings: settings,
		port:     port,
	}

	go controller.runWatchdogLoop()

	return controller, nil
}

func (c *Controller) runWatchdogLoop() {
	// @Todo: Do not leak this goroutine for 30 minutes.
	for {
		err := c.SendCommand(&WriteParameterRequest{
			ParameterID: NetParamWatchdogTTL,
			Parameter:   []byte{0x10, 0x0e, 0x00, 0x00}, // 3600 seconds = 1 hour
		})
		if err == io.EOF {
			return
		}
		if err != nil && c.settings.LogErrors {
			log.Println("[zigbee] failed to send watchdog reset command:", err)
		}
		time.Sleep(30 * time.Minute)
	}
}

func (c *Controller) Close() error {
	return c.port.Close()
}

func (c *Controller) Start() (chan zigbee.IncomingMessage, error) {
	messages := make(chan zigbee.IncomingMessage, 1)
	go func() {
		defer close(messages)
		for {
			data, err := slip.ReadPacket(c.port)
			if err == io.EOF {
				return
			}
			if err != nil {
				if c.settings.LogErrors {
					log.Println("[zigbee]", err)
					log.Println("[zigbee] exiting")
				}
				return
			}

			frame, err := ParseFrame(data, true)
			if err != nil {
				if c.settings.LogErrors {
					log.Println("[zigbee] failed to parse frame:", err)
				}
				continue
			}

			if c.settings.LogCommands {
				fmt.Printf("<-- %T%+v\n", frame.Command, frame.Command)
			}

			switch cmd := frame.Command.(type) {
			case *ReceivedDataNotification:
				c.SendCommand(&ReadReceivedDataRequest{})

			case *ReadReceivedDataResponse:
				messages <- zigbee.IncomingMessage{
					Source:              cmd.Source,
					SourceEndpoint:      cmd.SourceEndpoint,
					DestinationEndpoint: cmd.DestinationEndpoint,
					ClusterID:           cmd.ClusterID,
					LinkQuality:         cmd.LQI,
					Data:                cmd.Payload,
				}
			}
		}
	}()
	return messages, nil
}

func (c *Controller) Send(msg zigbee.OutgoingMessage) error {
	id := atomic.AddUint32(&c.requestSequence, 1)
	return c.SendCommand(&EnqueueSendDataRequest{
		RequestID:           uint8(id),
		Destination:         msg.Destination,
		DestinationEndpoint: msg.DestinationEndpoint,
		ProfileID:           0x0104, // @Todo: Hardcoded.
		ClusterID:           msg.ClusterID,
		SourceEndpoint:      msg.SourceEndpoint,
		Payload:             msg.Data,
		TxOptions:           0,
		Radius:              msg.Radius,
	})
}

func (c *Controller) PermitJoining(enabled bool) error {
	var value byte
	if enabled {
		value = 59
	}
	return c.SendCommand(&WriteParameterRequest{
		ParameterID: NetParam(0x21),
		Parameter:   []byte{value},
	})
}

func (c *Controller) SendCommand(command SerializableCommand) error {
	if c.settings.LogCommands {
		fmt.Printf("--> %T%+v\n", command, command)
	}

	sequence := atomic.AddUint32(&c.sequence, 1)
	data, err := SerializeFrame(Frame{
		SequenceNumber: uint8(sequence),
		Command:        command,
	})
	if err != nil {
		return err
	}

	return slip.WritePacket(c.port, data)
}
