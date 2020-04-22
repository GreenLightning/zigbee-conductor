package zigbee

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

type ControllerSettings struct {
	Port       string
	PermitJoin bool
}

type Controller struct {
	settings ControllerSettings
	port     *Port
}

func NewController(settings ControllerSettings) (*Controller, error) {
	callbacks := Callbacks{
		BeforeWrite: func(command interface{}) {
			fmt.Printf("--> %s%+v\n", reflect.TypeOf(command).Name(), command)
		},

		AfterRead: func(command interface{}) {
			fmt.Printf("<-- %s%+v\n", reflect.TypeOf(command).Name(), command)
		},

		OnReadError: func(err error) ErrorHandling {
			if errors.Is(err, ErrInvalidFrame) || errors.Is(err, ErrGarbage) {
				log.Println(err)
				return ErrorHandlingContinue
			}
			return ErrorHandlingPanic
		},

		OnParseError: func(err error, frame Frame) ErrorHandling {
			if err == ErrCommandInvalidFrame {
				log.Println("invalid frame")
				return ErrorHandlingContinue
			}
			if err == ErrCommandUnknownFrameHeader {
				log.Println("unknown frame:", frame)
				return ErrorHandlingContinue
			}
			return ErrorHandlingPanic
		},
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
		return err
	}

	_, err = c.port.WriteCommand(SysVersionRequest{})
	if err != nil {
		return err
	}

	response, err := c.port.WriteCommand(UtilGetDeviceInfoRequest{})
	if err != nil {
		return err
	}

	if response := response.(UtilGetDeviceInfoResponse); response.DeviceState != DeviceStateCoordinator {
		handler := c.port.RegisterOneOffHandler(ZdoStateChangeInd{})
		_, err = c.port.WriteCommand(ZdoStartupFromAppRequest{StartDelay: 100})
		if err != nil {
			return err
		}
		_, err := handler.Receive()
		if err != nil {
			return err
		}
	}

	handler := c.port.RegisterOneOffHandler(ZdoActiveEP{})
	_, err = c.port.WriteCommand(ZdoActiveEPRequest{})
	if err != nil {
		return err
	}
	cmd, err := handler.Receive()
	if err != nil {
		return err
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
				return err
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
		return err
	}

	return nil
}

func (c *Controller) Close() error {
	return c.port.Close()
}

func (c *Controller) RegisterPermanentHandler(commandPrototype interface{}) *Handler {
	return c.port.RegisterPermanentHandler(commandPrototype)
}
