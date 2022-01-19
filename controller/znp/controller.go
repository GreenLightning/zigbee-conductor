// Implements a Controller for Texas Instrument's CC253X-based dongles.
//
// For documentation of the serial protocol, please refer to the following
// documents: "ZNP Interface Specification", "Z-Stack Monitor and Test API".
package znp

import (
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/GreenLightning/zigbee-conductor/zigbee"
)

type Controller struct {
	settings zigbee.ControllerSettings
	sequence uint32
	port     *Port
}

func NewController(settings zigbee.ControllerSettings) (*Controller, error) {
	callbacks := Callbacks{
		OnReadError: func(err error) ErrorHandling {
			if errors.Is(err, ErrInvalidFrame) || errors.Is(err, ErrGarbage) {
				if settings.LogErrors {
					log.Println("[zigbee]", err)
				}
				return ErrorHandlingContinue
			}
			return ErrorHandlingPanic
		},

		OnParseError: func(err error, frame Frame) ErrorHandling {
			if err == ErrCommandInvalidFrame {
				if settings.LogErrors {
					log.Println("[zigbee] invalid serial frame")
				}
				return ErrorHandlingContinue
			}
			if err == ErrCommandUnknownFrameHeader {
				if settings.LogErrors {
					log.Println("[zigbee] unknown serial frame:", frame)
				}
				return ErrorHandlingContinue
			}
			return ErrorHandlingPanic
		},
	}

	if settings.LogCommands {
		callbacks.BeforeWrite = func(command interface{}) {
			fmt.Printf("--> %T%+v\n", command, command)
		}
		callbacks.AfterRead = func(command interface{}) {
			fmt.Printf("<-- %T%+v\n", command, command)
		}
	}

	port, err := NewPort(settings.Port, callbacks)
	if err != nil {
		return nil, err
	}

	return &Controller{
		settings: settings,
		port:     port,
	}, nil
}

type Endpoint struct {
	Endpoint  uint8
	AppProfID uint16
}

var endpoints = []Endpoint{
	Endpoint{1, 0x0104},
	Endpoint{2, 0x0101},
	Endpoint{3, 0x0105},
	Endpoint{4, 0x0107},
	Endpoint{5, 0x0108},
	Endpoint{6, 0x0109},
	Endpoint{8, 0x0104},
}

func (c *Controller) Start() (chan zigbee.IncomingMessage, error) {
	err := c.port.WriteMagicByteForBootloader()
	if err != nil {
		return nil, fmt.Errorf("writing magic byte for bootloader: %w", err)
	}

	// This command is used to test communication with the device.
	// The response may be slow if the device has to finish booting,
	// therefore we use a custom timeout value.
	_, err = c.port.WriteCommandTimeout(SysVersionRequest{}, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("getting system version: %w", err)
	}

	response, err := c.port.WriteCommand(UtilGetDeviceInfoRequest{})
	if err != nil {
		return nil, fmt.Errorf("getting device info: %w", err)
	}

	if response := response.(UtilGetDeviceInfoResponse); response.DeviceState != DeviceStateCoordinator {
		handler := c.port.RegisterOneOffHandler(ZdoStateChangeInd{})
		_, err = c.port.WriteCommand(ZdoStartupFromAppRequest{StartDelay: 100})
		if err != nil {
			return nil, fmt.Errorf("sending startup from app: %w", err)
		}
		_, err := handler.Receive()
		if err != nil {
			return nil, fmt.Errorf("waiting for state change: %w", err)
		}
	}

	// Activate endpoints to receive incoming messages.
	for _, endpoint := range endpoints {
		_, err = c.port.WriteCommand(AfRegisterRequest{
			Endpoint:    endpoint.Endpoint,
			AppProfID:   endpoint.AppProfID,
			AppDeviceID: 0x0005,
			LatencyReq:  LatencyReqNoLatency,
		})
		if err != nil {
			return nil, fmt.Errorf("sending register endpoint: %w", err)
		}
	}

	handler := c.RegisterPermanentHandler(AfIncomingMsg{})
	output := make(chan zigbee.IncomingMessage)

	go func() {
		for {
			cmd, err := handler.Receive()
			if err != nil {
				// @Todo: Forward error to user.
				break
			}

			message := cmd.(AfIncomingMsg)
			output <- zigbee.IncomingMessage{
				Source: zigbee.Address{
					Mode:  zigbee.AddressModeNWK,
					Short: message.SrcAddr,
				},
				SourceEndpoint:      message.SrcEndpoint,
				DestinationEndpoint: message.DstEndpoint,
				ClusterID:           message.ClusterID,
				LinkQuality:         message.LinkQuality,
				Data:                message.Data,
			}
		}
		close(output)
	}()

	return output, nil
}

func (c *Controller) Close() error {
	return c.port.Close()
}

func (c *Controller) Send(message zigbee.OutgoingMessage) error {
	if mode := message.Destination.Mode; mode != zigbee.AddressModeNWK && mode != zigbee.AddressModeCombined {
		return fmt.Errorf("address mode not supported: %v", mode)
	}

	sequence := atomic.AddUint32(&c.sequence, 1)

	_, err := c.WriteCommand(AfDataRequest{
		DstAddr:        message.Destination.Short,
		DstEndpoint:    message.DestinationEndpoint,
		SrcEndpoint:    message.SourceEndpoint,
		ClusterID:      message.ClusterID,
		TransSeqNumber: uint8(sequence),
		Options:        0,
		Radius:         message.Radius,
		Data:           message.Data,
	})

	return err
}

func (c *Controller) PermitJoining(enabled bool) error {
	permitJoinRequest := ZdoMgmtPermitJoinRequest{
		AddrMode: 0x0f,
		DstAddr:  0xfffc, // broadcast
	}
	if enabled {
		permitJoinRequest.Duration = 0xff // turn on indefinitely
	}
	_, err := c.port.WriteCommand(permitJoinRequest)
	if err != nil {
		return fmt.Errorf("sending permit join: %w", err)
	}

	return nil
}

func (c *Controller) RegisterPermanentHandler(commandPrototype interface{}) *Handler {
	return c.port.RegisterPermanentHandler(commandPrototype)
}

func (c *Controller) WriteCommand(command interface{}) (interface{}, error) {
	return c.port.WriteCommand(command)
}

func (c *Controller) WriteCommandTimeout(command interface{}, timeout time.Duration) (interface{}, error) {
	return c.port.WriteCommandTimeout(command, timeout)
}
