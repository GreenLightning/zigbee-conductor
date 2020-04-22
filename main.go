package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"reflect"

	"github.com/GreenLightning/zigbee-conductor/zigbee"
)

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

func main() {
	portFlag := flag.String("port", "/dev/ttyACM0", "name of the serial port to use")
	permitJoinFlag := flag.Bool("permitJoin", false, "permit devices to join the network")

	flag.Parse()

	callbacks := zigbee.Callbacks{
		BeforeWrite: func(command interface{}) {
			fmt.Printf("--> %s%+v\n", reflect.TypeOf(command).Name(), command)
		},

		AfterRead: func(command interface{}) {
			fmt.Printf("<-- %s%+v\n", reflect.TypeOf(command).Name(), command)
		},

		OnReadError: func(err error) zigbee.ErrorHandling {
			if errors.Is(err, zigbee.ErrInvalidFrame) || errors.Is(err, zigbee.ErrGarbage) {
				log.Println(err)
				return zigbee.ErrorHandlingContinue
			}
			return zigbee.ErrorHandlingPanic
		},

		OnParseError: func(err error, frame zigbee.Frame) zigbee.ErrorHandling {
			if err == zigbee.ErrCommandInvalidFrame {
				log.Println("invalid frame")
				return zigbee.ErrorHandlingContinue
			}
			if err == zigbee.ErrCommandUnknownFrameHeader {
				log.Println("unknown frame:", frame)
				return zigbee.ErrorHandlingContinue
			}
			return zigbee.ErrorHandlingPanic
		},
	}

	port, err := zigbee.NewPort(*portFlag, callbacks)
	check(err)

	defer port.Close()

	check(port.WriteMagicByteForBootloader())

	_, err = port.WriteCommand(zigbee.SysVersionRequest{})
	check(err)

	response, err := port.WriteCommand(zigbee.UtilGetDeviceInfoRequest{})
	check(err)

	if response := response.(zigbee.UtilGetDeviceInfoResponse); response.DeviceState != zigbee.DeviceStateCoordinator {
		handler := port.RegisterHandler(zigbee.ZdoStateChangeInd{})
		_, err = port.WriteCommand(zigbee.ZdoStartupFromAppRequest{StartDelay: 100})
		check(err)
		_, err := handler.Receive()
		check(err)
	}

	handler := port.RegisterHandler(zigbee.ZdoActiveEP{})
	_, err = port.WriteCommand(zigbee.ZdoActiveEPRequest{})
	check(err)
	cmd, err := handler.Receive()
	check(err)

	activeEPs := cmd.(zigbee.ZdoActiveEP)
	for _, endpoint := range endpoints {
		found := false
		for _, ep := range activeEPs.ActiveEPs {
			if ep == endpoint.Endpoint {
				found = true
				break
			}
		}

		if !found {
			_, err = port.WriteCommand(zigbee.AfRegisterRequest{
				Endpoint:    endpoint.Endpoint,
				AppProfID:   endpoint.AppProfID,
				AppDeviceID: 0x0005,
				LatencyReq:  zigbee.LatencyReqNoLatency,
			})
		}
	}

	permitJoinRequest := zigbee.ZdoMgmtPermitJoinRequest{
		AddrMode: 0x0f,
		DstAddr:  0xfffc, // broadcast
	}
	if *permitJoinFlag {
		permitJoinRequest.Duration = 0xff // turn on indefinitely
	}
	_, err = port.WriteCommand(permitJoinRequest)
	check(err)

	select {}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
