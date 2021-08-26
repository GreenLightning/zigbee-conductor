package conbee

import (
	"fmt"
	"io"

	"github.com/GreenLightning/zigbee-conductor/pkg/slip"
	"github.com/GreenLightning/zigbee-conductor/zigbee"
	"github.com/jacobsa/go-serial/serial"
)

type Controller struct {
	settings zigbee.ControllerSettings
	port     io.ReadWriteCloser
}

// @Todo: Add logging.
func NewController(settings zigbee.ControllerSettings) (zigbee.Controller, error) {
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

	return &Controller{
		settings: settings,
		port:     port,
	}, nil
}

func (c *Controller) Close() error {
	return c.port.Close()
}

func (c *Controller) Start() (chan zigbee.IncomingMessage, error) {
	messages := make(chan zigbee.IncomingMessage, 1)
	go func() {
		for {
			data, err := slip.ReadPacket(c.port)
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}

			frame, err := ParseFrame(data, true)
			if err != nil {
				continue
			}

			// @Todo: Handle frame.
			fmt.Println(frame)
		}
	}()
	return messages, nil
}

func (c *Controller) Send(msg zigbee.OutgoingMessage) error {
	panic("not implemented")
}

func (c *Controller) PermitJoining(enabled bool) error {
	panic("not implemented")
}
