# github.com/rickbassham/gpsd

A simple go client library for interfacing with gpsd.

## Example

```go
package main

import (
	"fmt"
	"net"

	"github.com/rickbassham/gpsd"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.8:2947")
	if err != nil {
		panic(err.Error())
	}

	svc := gpsd.New(conn)

	svc.Watch()

	messageTypes := gpsd.MessageTypeTimePositionVelocity | gpsd.MessageTypeVersion

	messages, errors := svc.Scan(messageTypes)

	// Handler for TPV reports
	tpv := func(rpt gpsd.TimePositionVelocityReport) {
		fmt.Println(fmt.Sprintf("Location: %0.7f, %0.7f", rpt.Lat, rpt.Lon))

		// only care about the first location event, so stop reading from gpsd when we get it
		svc.Stop()
	}

	// Handler for VERSION reports
	ver := func(rpt gpsd.VersionReport) {
		fmt.Println(fmt.Sprintf("Connected to GPSD version %s", rpt.Release))
	}

	h := &gpsd.Handler{}
	h = h.WithTimePositionVelocity(tpv).WithVersion(ver)

	gpsd.ScanLoop(messages, errors, h, func(err error) {
		panic(err.Error())
	})

	// Output:
	// Connected to GPSD version 3.17
	// Location: 40.7828687, -73.9675438
}
```