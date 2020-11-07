package zigbee

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type ControllerSettings struct {
	Port        string
	PermitJoin  bool
	LogCommands bool
	LogErrors   bool
}

type Controller struct {
	settings ControllerSettings
	port     *Port
}

func NewController(settings ControllerSettings) (*Controller, error) {
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
					log.Println("[zigbee] invalid frame")
				}
				return ErrorHandlingContinue
			}
			if err == ErrCommandUnknownFrameHeader {
				if settings.LogErrors {
					log.Println("[zigbee] unknown frame:", frame)
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

func (c *Controller) Start() error {
	err := c.port.WriteMagicByteForBootloader()
	if err != nil {
		return fmt.Errorf("writing magic byte for bootloader: %w", err)
	}

	// This command is used to test communication with the device.
	// The response may be slow if the device has to finish booting,
	// therefore we use a custom timeout value.
	_, err = c.port.WriteCommandTimeout(SysVersionRequest{}, 10*time.Second)
	if err != nil {
		return fmt.Errorf("getting system version: %w", err)
	}

	response, err := c.port.WriteCommand(UtilGetDeviceInfoRequest{})
	if err != nil {
		return fmt.Errorf("getting device info: %w", err)
	}

	if response := response.(UtilGetDeviceInfoResponse); response.DeviceState != DeviceStateCoordinator {
		handler := c.port.RegisterOneOffHandler(ZdoStateChangeInd{})
		_, err = c.port.WriteCommand(ZdoStartupFromAppRequest{StartDelay: 100})
		if err != nil {
			return fmt.Errorf("sending startup from app: %w", err)
		}
		_, err := handler.Receive()
		if err != nil {
			return fmt.Errorf("waiting for state change: %w", err)
		}
	}

	handler := c.port.RegisterOneOffHandler(ZdoActiveEP{})
	_, err = c.port.WriteCommand(ZdoActiveEPRequest{})
	if err != nil {
		return fmt.Errorf("getting active endpoints: %w", err)
	}
	cmd, err := handler.Receive()
	if err != nil {
		return fmt.Errorf("waiting for active endpoints: %w", err)
	}

	activeEPs := cmd.(ZdoActiveEP)
	for _, endpoint := range endpoints {
		found := false
		for _, ep := range activeEPs.ActiveEPs {
			if ep == endpoint.Endpoint {
				found = true
				break
			}
		}

		if !found {
			_, err = c.port.WriteCommand(AfRegisterRequest{
				Endpoint:    endpoint.Endpoint,
				AppProfID:   endpoint.AppProfID,
				AppDeviceID: 0x0005,
				LatencyReq:  LatencyReqNoLatency,
			})
			if err != nil {
				return fmt.Errorf("sending register endpoint: %w", err)
			}
		}
	}

	permitJoinRequest := ZdoMgmtPermitJoinRequest{
		AddrMode: 0x0f,
		DstAddr:  0xfffc, // broadcast
	}
	if c.settings.PermitJoin {
		permitJoinRequest.Duration = 0xff // turn on indefinitely
	}
	_, err = c.port.WriteCommand(permitJoinRequest)
	if err != nil {
		return fmt.Errorf("sending permit join: %w", err)
	}

	return nil
}

func (c *Controller) Close() error {
	return c.port.Close()
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
