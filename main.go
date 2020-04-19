package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"reflect"
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

	port, err := NewPort(*portFlag, callbacks)
	check(err)

	defer port.Close()

	check(port.WriteMagicByteForBootloader())

	_, err = port.WriteCommand(SysVersionRequest{})
	check(err)

	response, err := port.WriteCommand(UtilGetDeviceInfoRequest{})
	check(err)

	if response := response.(UtilGetDeviceInfoResponse); response.DeviceState != DeviceStateCoordinator {
		handler := port.RegisterHandler(ZdoStateChangeInd{})
		_, err = port.WriteCommand(ZdoStartupFromAppRequest{StartDelay: 100})
		check(err)
		_, err := handler.Receive()
		check(err)
	}

	handler := port.RegisterHandler(ZdoActiveEP{})
	_, err = port.WriteCommand(ZdoActiveEPRequest{})
	check(err)
	cmd, err := handler.Receive()
	check(err)

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
			_, err = port.WriteCommand(AfRegisterRequest{
				Endpoint:    endpoint.Endpoint,
				AppProfID:   endpoint.AppProfID,
				AppDeviceID: 0x0005,
				LatencyReq:  LatencyReqNoLatency,
			})
		}
	}

	permitJoinRequest := ZdoMgmtPermitJoinRequest{
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
