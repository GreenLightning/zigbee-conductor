# Zigbee Conductor

This module allows interacting with the ZigBee network using USB dongles.
Currently supported are Texas Instrument's CC253X-based dongles and the ConBee II gateway from Phoscon.

While ZigBee is standardized, the different dongles expose their functionality
over a serial port connection using individual APIs. Therefore the `zigbee`
package provides a general `Controller` interface, which is implemented for the
different vendors by individual subpackages of the `controller` package:

- `controller/conbee` for the [ConBee II](https://phoscon.de/en/conbee2).
- `controller/znp` for CC253X-based dongles (Zigbee Network Processor is the name of Texas Instrument's software).

The `controllerregistry` package can be used to dynamically create a controller
from a string identifier (and to register external controller implementations
for other USB dongles).

A `Controller` handles most of the ZigBee stack including the Physical (PHY),
MAC, Network (NWK) and Application Support (APS) layers. For higher-level
support the `zcl` package provides functions to parse and serialize frames of
the ZigBee Cluster Library, which sits above the APS layer.

# Documentation

Many of the constants and structs are transcribed from the specification.
Please refer to the following documents for more information.

- The Zigbee and Zigbee Cluster Library Specifications
  (available for download from [here](https://zigbeealliance.org/solution/zigbee/)).
- The book "Zigbee Wireless Networking" by Drew Gislason for a general overview
  and explanation of Zigbee concepts.
- Specifications for the serial protocols used by the dongles linked in the
  package documentation of the controller packages.

# Project Status

The code is robust, but incomplete. It has been running in my home automation
setup for years, however my setup is relatively simple and ZigBee is only used
to read sensor values and receive events from buttons. I am in the process of
integrating the library into another project, where it will be responsible for
commissioning and controlling devices as well.

There might be breaking changes as I work on the API and integrate new features.
Due to the size of the protocol, things are implemented on an as-needed basis.
Feel free to open an issue or pull request if you need a specific command.
