package main

import (
	"flag"
	"fmt"

	"github.com/GreenLightning/zigbee-conductor/zcl"
	"github.com/GreenLightning/zigbee-conductor/zigbee"
)

func main() {
	portFlag := flag.String("port", "/dev/ttyACM0", "name of the serial port to use")
	permitJoinFlag := flag.Bool("permitJoin", false, "permit devices to join the network")

	flag.Parse()

	controller, err := zigbee.NewController(zigbee.ControllerSettings{
		Port:        *portFlag,
		PermitJoin:  *permitJoinFlag,
		LogCommands: true,
		LogErrors:   true,
	})
	check(err)

	defer controller.Close()

	handler := controller.RegisterPermanentHandler(zigbee.AfIncomingMsg{})

	check(controller.Start())

	for {
		cmd, err := handler.Receive()
		check(err)
		msg := cmd.(zigbee.AfIncomingMsg)
		frame, err := zcl.ParseFrame(msg.ClusterID, msg.Data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("<--- %T%+v\n", frame, frame)
		if frame.Type == zcl.FRAME_TYPE_GLOBAL && !frame.ManufacturerSpecific && frame.CommandIdentifier == zcl.COMMAND_ID_REPORT_ATTRIBUTES {
			reportCmd, err := zcl.ParseReportAttributesCommand(frame.Data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, report := range reportCmd.Reports {
				fmt.Printf("<---- %T%+v\n", report, report)
			}
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
