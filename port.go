package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type ErrorHandling int

const (
	ErrorHandlingPanic    ErrorHandling = 0
	ErrorHandlingStop     ErrorHandling = 1
	ErrorHandlingContinue ErrorHandling = 1
)

type Callbacks struct {
	BeforeWrite func(command interface{})
	AfterRead   func(command interface{})

	OnCommand func(command interface{})

	OnReadError  func(err error) ErrorHandling
	OnParseError func(err error, frame Frame) ErrorHandling
}

type responseHandler chan interface{}

type Port struct {
	sp  io.ReadWriteCloser
	cbs Callbacks

	handlerMutex sync.Mutex
	handlers     map[FrameHeader]responseHandler
}

func NewPort(name string, callbacks Callbacks) (*Port, error) {
	options := serial.OpenOptions{
		PortName:          name,
		BaudRate:          115200,
		RTSCTSFlowControl: true,
		DataBits:          8,
		StopBits:          1,
		MinimumReadSize:   1,
	}

	sp, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	port := &Port{
		sp:  sp,
		cbs: callbacks,

		handlerMutex: sync.Mutex{},
		handlers:     make(map[FrameHeader]responseHandler),
	}

	go port.loop()

	return port, nil
}

func (p *Port) Close() error {
	return p.sp.Close()
}

func (p *Port) WriteMagicByteForBootloader() error {
	// Send magic byte to bootloader, to skip 60s wait after startup.
	//
	// "I found the solution.
	//  The serial boot loader is waiting for 60 seconds before jumping to the ZNP application image, as documented in Serial Boot Loader document.
	//  Sending a magic byte 0xef can force the sbl to skip the waiting.
	//  However, the version in Home-1.2.1 changed the magic byte from 0x10 to 0xef, and the document did not update this.
	//  That's why I could not get it work before.
	//  Also the response to the magic byte has been changed from 0x00 to SYS_RESET_IND. again no document update."
	//
	//    - Kyle Zhou in [2]
	//
	// See: [1] https://github.com/Koenkk/zigbee2mqtt/issues/1343
	// See: [2] https://e2e.ti.com/support/wireless-connectivity/zigbee-and-thread/f/158/p/160948/1361000#1361000
	//
	_, err := p.sp.Write([]byte{0xEF})
	return err
}

func (p *Port) WriteCommand(command interface{}) (interface{}, error) {
	if p.cbs.BeforeWrite != nil {
		p.cbs.BeforeWrite(command)
	}

	frame := buildFrameForCommand(command)

	var handler responseHandler
	var responseHeader FrameHeader

	if frame.Type == FRAME_TYPE_SREQ {
		handler = make(responseHandler)

		responseHeader = frame.FrameHeader
		responseHeader.Type = FRAME_TYPE_SRSP

		p.handlerMutex.Lock()
		p.handlers[responseHeader] = handler
		p.handlerMutex.Unlock()
	}

	err := writeFrame(p.sp, frame)
	if err != nil {
		return nil, err
	}

	if handler == nil {
		return nil, nil
	}

	var response interface{}
	select {
	case response = <-handler:
	case <-time.After(1000 * time.Millisecond):
	}

	p.handlerMutex.Lock()
	delete(p.handlers, responseHeader)
	p.handlerMutex.Unlock()

	if response == nil {
		select {
		case response = <-handler:
		default:
		}
	}

	if response == nil {
		return nil, errors.New("response timeout")
	}

	return response, nil
}

func (p *Port) loop() {
	r := bufio.NewReaderSize(p.sp, 256)
	for {
		frame, err := readFrame(r)
		if err != nil {
			if errors.Is(err, os.ErrClosed) {
				break
			}
			var handling ErrorHandling
			if p.cbs.OnReadError != nil {
				handling = p.cbs.OnReadError(err)
			}
			if handling == ErrorHandlingContinue {
				continue
			} else if handling == ErrorHandlingStop {
				break
			} else {
				panic(err)
			}
		}

		command, err := parseCommandFromFrame(frame)
		if err != nil {
			handling := ErrorHandlingStop
			if p.cbs.OnParseError != nil {
				handling = p.cbs.OnParseError(err, frame)
			}
			if handling == ErrorHandlingContinue {
				continue
			} else if handling == ErrorHandlingStop {
				break
			} else {
				panic(err)
			}
		}

		if p.cbs.AfterRead != nil {
			p.cbs.AfterRead(command)
		}

		var handler responseHandler
		p.handlerMutex.Lock()
		handler = p.handlers[frame.FrameHeader]
		p.handlerMutex.Unlock()

		if handler != nil {
			handler <- command
		} else if p.cbs.OnCommand != nil {
			p.cbs.OnCommand(command)
		}
	}
}
