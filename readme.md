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
