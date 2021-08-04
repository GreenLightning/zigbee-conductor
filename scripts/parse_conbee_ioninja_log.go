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

		fmt.Printf("%s %s %+v %T%+v\n", p.Timestamp, p.Dir, frame, frame.Command, frame.Command)
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
