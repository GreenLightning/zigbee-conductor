package main

import (
	"flag"
	"fmt"

	"github.com/GreenLightning/zigbee-conductor/controller/controllerregistry"
	"github.com/GreenLightning/zigbee-conductor/zcl"
	"github.com/GreenLightning/zigbee-conductor/zigbee"
)

func main() {
	portFlag := flag.String("port", "/dev/ttyACM0", "name of the serial port to use")
	controllerFlag := flag.String("controller", "znp", "type of the controller: conbee, znp")
	permitJoinFlag := flag.Bool("permitJoin", false, "permit devices to join the network")

	flag.Parse()

	controller, err := controllerregistry.NewController(*controllerFlag, zigbee.ControllerSettings{
		Port:        *portFlag,
		LogCommands: true,
		LogErrors:   true,
	})
	check(err)

	defer controller.Close()

	incoming, err := controller.Start()
	check(err)

	controller.PermitJoining(*permitJoinFlag)

	for message := range incoming {
		frame, err := zcl.ParseFrame(message.Data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("<--- %T%+v\n", frame, frame)
		if frame.Type == zcl.FrameTypeGlobal && !frame.ManufacturerSpecific {
			if frame.CommandID == zcl.CommandReadAttributes {
				cmd2, err := zcl.ParseReadAttributesCommand(frame.Data)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, attribute := range cmd2.Attributes {
					fmt.Printf("<---- %v\n", attribute)
				}

			} else if frame.CommandID == zcl.CommandReadAttributesResponse {
				cmd2, err := zcl.ParseReadAttributesResponseCommand(frame.Data)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, record := range cmd2.Records {
					fmt.Printf("<---- %T%+v\n", record, record)
				}

			} else if frame.CommandID == zcl.CommandReportAttributes {
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
