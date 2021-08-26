package controllerregistry

import (
	"errors"
	"fmt"

	"github.com/GreenLightning/zigbee-conductor/zigbee"
)

type ControllerFactory func(settings zigbee.ControllerSettings) (zigbee.Controller, error)

var ErrNotFound = errors.New("controller not found")

var registry = make(map[string]ControllerFactory)

func Register(name string, factory ControllerFactory) {
	if _, ok := registry[name]; ok {
		panic(fmt.Sprintf("controller already registered: %s", name))
	}
	registry[name] = factory
}

func NewController(name string, settings zigbee.ControllerSettings) (zigbee.Controller, error) {
	factory, ok := registry[name]
	if !ok {
		return nil, ErrNotFound
	}

	return factory(settings)
}
