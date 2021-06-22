package znp

import (
	"bufio"
	"errors"
	"fmt"
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
	ErrorHandlingContinue ErrorHandling = 2
)

type Callbacks struct {
	BeforeWrite func(command interface{})
	AfterRead   func(command interface{})

	OnReadError  func(err error) ErrorHandling
	OnParseError func(err error, frame Frame) ErrorHandling
}

var ErrTimeout = errors.New("command timed out")

type handlerResult struct {
	command interface{}
	err     error
}

type Handler struct {
	results chan handlerResult
	timer   *time.Timer
}

func newHandler() *Handler {
	return &Handler{
		results: make(chan handlerResult, 1),
	}
}

func (h *Handler) fail() {
	h.results <- handlerResult{nil, ErrTimeout}
}

func (h *Handler) fulfill(command interface{}) {
	h.results <- handlerResult{command, nil}
	if h.timer != nil {
		h.timer.Stop()
	}
}

func (h *Handler) Receive() (interface{}, error) {
	result := <-h.results
	return result.command, result.err
}

type Port struct {
	sp  io.ReadWriteCloser
	cbs Callbacks

	handlerMutex sync.Mutex
	handlers     map[FrameHeader]*Handler
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
		handlers:     make(map[FrameHeader]*Handler),
	}

	go port.loop()

	return port, nil
}

func (p *Port) Close() error {
	return p.sp.Close()
}

func (p *Port) RegisterOneOffHandler(commandPrototype interface{}) *Handler {
	header := getHeaderForCommand(commandPrototype)
	return p.registerHandler(header, 10*time.Second)
}

func (p *Port) RegisterPermanentHandler(commandPrototype interface{}) *Handler {
	header := getHeaderForCommand(commandPrototype)
	return p.registerHandler(header, 0)
}

func (p *Port) registerHandler(header FrameHeader, timeout time.Duration) *Handler {
	handler := newHandler()

	p.handlerMutex.Lock()
	defer p.handlerMutex.Unlock()

	if _, ok := p.handlers[header]; ok {
		panic(fmt.Sprintf("handler for %v already exists", header))
	}
	p.handlers[header] = handler

	if timeout != 0 {
		handler.timer = time.AfterFunc(timeout, func() {
			p.removeHandler(handler, header)
		})
	}

	return handler
}

func (p *Port) removeHandler(handler *Handler, header FrameHeader) {
	found := false
	p.handlerMutex.Lock()
	if p.handlers[header] == handler {
		delete(p.handlers, header)
		found = true
	}
	p.handlerMutex.Unlock()

	if found {
		handler.fail()
	}
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
	return p.WriteCommandTimeout(command, 1*time.Second)
}

func (p *Port) WriteCommandTimeout(command interface{}, timeout time.Duration) (interface{}, error) {
	if p.cbs.BeforeWrite != nil {
		p.cbs.BeforeWrite(command)
	}

	frame := buildFrameForCommand(command)

	var handler *Handler

	if frame.Type == FRAME_TYPE_SREQ {
		responseHeader := frame.FrameHeader
		responseHeader.Type = FRAME_TYPE_SRSP
		handler = p.registerHandler(responseHeader, timeout)
	}

	err := writeFrame(p.sp, frame)
	if err != nil {
		return nil, err
	}

	if handler != nil {
		return handler.Receive()
	}

	return nil, nil
}

func (p *Port) loop() {
	r := bufio.NewReaderSize(p.sp, 256)
	for {
		frame, err := readFrame(r)
		if err == io.ErrNoProgress {
			continue
		}
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

		p.handlerMutex.Lock()
		handler := p.handlers[frame.FrameHeader]
		if handler != nil && handler.timer != nil {
			delete(p.handlers, frame.FrameHeader)
		}
		p.handlerMutex.Unlock()

		if handler != nil {
			handler.fulfill(command)
		}
	}
}
