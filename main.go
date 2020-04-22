package main

import (
	"flag"

	"github.com/GreenLightning/zigbee-conductor/zigbee"
)

func main() {
	portFlag := flag.String("port", "/dev/ttyACM0", "name of the serial port to use")
	permitJoinFlag := flag.Bool("permitJoin", false, "permit devices to join the network")

	flag.Parse()

	controller, err := zigbee.NewController(zigbee.ControllerSettings{
		Port:       *portFlag,
		PermitJoin: *permitJoinFlag,
	})
	check(err)

	defer controller.Close()

	check(controller.Start())

	select {}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
