package main

import (
	"flag"
	"fmt"

	"github.com/GreenLightning/zigbee-conductor/zcl"
	"github.com/GreenLightning/zigbee-conductor/znp"
)

func main() {
	portFlag := flag.String("port", "/dev/ttyACM0", "name of the serial port to use")
	permitJoinFlag := flag.Bool("permitJoin", false, "permit devices to join the network")

	flag.Parse()

	controller, err := znp.NewController(znp.ControllerSettings{
		Port:        *portFlag,
		PermitJoin:  *permitJoinFlag,
		LogCommands: true,
		LogErrors:   true,
	})
	check(err)

	defer controller.Close()

	handler := controller.RegisterPermanentHandler(znp.AfIncomingMsg{})

	check(controller.Start())

	for {
		cmd, err := handler.Receive()
		check(err)
		msg := cmd.(znp.AfIncomingMsg)
		frame, err := zcl.ParseFrame(msg.Data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("<--- %T%+v\n", frame, frame)
		if frame.Type == zcl.FrameTypeGlobal && !frame.ManufacturerSpecific {
			if frame.CommandID == zcl.CommandReportAttributes {
				cmd2, err := zcl.ParseReportAttributesCommand(frame.Data)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, report := range cmd2.Reports {
					fmt.Printf("<---- %T%+v\n", report, report)
				}
			} else if frame.CommandID == zcl.CommandReadReportingConfigurationResponse {
				cmd2, err := zcl.ParseReadReportingConfigurationResponseCommand(frame.Data)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, record := range cmd2.Records {
					fmt.Printf("<---- %T%+v\n", record, record)
				}
			}
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
