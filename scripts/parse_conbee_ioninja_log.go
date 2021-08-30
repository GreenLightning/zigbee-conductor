package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/GreenLightning/zigbee-conductor/conbee"
	"github.com/GreenLightning/zigbee-conductor/pkg/slip"
	"github.com/GreenLightning/zigbee-conductor/zcl"
)

type Parser struct {
	Timestamp string
	Dir       string
	Data      []byte
}

func (p *Parser) ParseLine(line string) {
	if len(line) < 75 {
		p.Flush()
		fmt.Println(line)
		return
	}

	timestamp := line[:19]
	dir := line[20:21]

	if dir != ">" && dir != "<" {
		p.Flush()
		fmt.Println(line)
		return
	}

	data, err := hex.DecodeString(strings.ReplaceAll(line[28:75], " ", ""))
	if err != nil {
		panic(err)
	}

	if p.Timestamp != "" && dir == p.Dir {
		p.Data = append(p.Data, data...)
		return
	}

	p.Flush()
	p.Timestamp = timestamp
	p.Dir = dir
	p.Data = data
}

func (p *Parser) Flush() {
	if p.Timestamp == "" {
		return
	}

	reader := bytes.NewReader(p.Data)
	for {
		packet, err := slip.ReadPacket(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Println(p.Timestamp, p.Dir, err)
			}
			break
		}

		incoming := p.Dir == ">"
		frame, err := conbee.ParseFrame(packet, incoming)
		if err != nil {
			fmt.Println(p.Timestamp, p.Dir, err)
			continue
		}

		spacer := strings.Repeat(" ", len(p.Timestamp))
		fmt.Printf("%s %s %+v\n", p.Timestamp, p.Dir, frame)
		fmt.Printf("%s %s   %T%+v\n", spacer, p.Dir, frame.Command, frame.Command)

		var payload []byte
		if request, ok := frame.Command.(*conbee.EnqueueSendDataRequest); ok {
			payload = request.Payload
		} else if response, ok := frame.Command.(*conbee.ReadReceivedDataResponse); ok {
			payload = response.Payload
		}

		if payload != nil {
			frame, err := zcl.ParseFrame(payload)
			if err != nil {
				continue
			}

			fmt.Printf("%s %s     %T%+v\n", spacer, p.Dir, frame, frame)
			if frame.Type == zcl.FrameTypeGlobal && !frame.ManufacturerSpecific {
				if frame.CommandID == zcl.CommandReadAttributes {
					cmd2, err := zcl.ParseReadAttributesCommand(frame.Data)
					if err != nil {
						fmt.Printf("%s %s       %v\n", spacer, p.Dir, err)
						continue
					}
					for _, attribute := range cmd2.Attributes {
						fmt.Printf("%s %s       %v\n", spacer, p.Dir, attribute)
					}

				} else if frame.CommandID == zcl.CommandReadAttributesResponse {
					cmd2, err := zcl.ParseReadAttributesResponseCommand(frame.Data)
					if err != nil {
						fmt.Printf("%s %s       %v\n", spacer, p.Dir, err)
						continue
					}
					for _, record := range cmd2.Records {
						fmt.Printf("%s %s       %T%+v\n", spacer, p.Dir, record, record)
					}

				} else if frame.CommandID == zcl.CommandReportAttributes {
					cmd2, err := zcl.ParseReportAttributesCommand(frame.Data)
					if err != nil {
						fmt.Printf("%s %s       %v\n", spacer, p.Dir, err)
						continue
					}
					for _, report := range cmd2.Reports {
						fmt.Printf("%s %s       %T%+v\n", spacer, p.Dir, report, report)
					}

				} else if frame.CommandID == zcl.CommandReadReportingConfigurationResponse {
					cmd2, err := zcl.ParseReadReportingConfigurationResponseCommand(frame.Data)
					if err != nil {
						fmt.Printf("%s %s       %v\n", spacer, p.Dir, err)
						continue
					}
					for _, record := range cmd2.Records {
						fmt.Printf("%s %s       %T%+v\n", spacer, p.Dir, record, record)
					}
				}
			}
		}
	}

	p.Timestamp = ""
	p.Dir = ""
	p.Data = nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: <filename>")
		os.Exit(2)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	var parser Parser
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parser.ParseLine(scanner.Text())
	}
	parser.Flush()

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
