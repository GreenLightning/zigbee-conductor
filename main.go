package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"reflect"
	"time"
)

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
		_, err = port.WriteCommand(ZdoStartupFromAppRequest{StartDelay: 100})
		check(err)
	}

	time.Sleep(1 * time.Second)

	permitJoinRequest := ZdoMgmtPermitJoinRequest{
		AddrMode: 0x0f,
		DstAddr:  0xfffc, // broadcast
	}
	if *permitJoinFlag {
		permitJoinRequest.Duration = 0xff // turn on indefinitely
	}
	_, err = port.WriteCommand(permitJoinRequest)
	check(err)

	time.Sleep(30 * time.Second)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
