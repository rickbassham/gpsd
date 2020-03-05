package gpsd_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"github.com/rickbassham/gpsd"
)

type MockConn struct {
	readableData *bytes.Buffer
	writtenData  *bytes.Buffer
}

func (c *MockConn) Read(p []byte) (n int, err error) {
	return c.readableData.Read(p)
}

func (c *MockConn) Write(p []byte) (n int, err error) {
	return c.writtenData.Write(p)
}

func TestGPSD(t *testing.T) {
	readable := bytes.Buffer{}

	readable.WriteString(`{"class":"VERSION","release":"3.17","rev":"3.17","proto_major":3,"proto_minor":12}`)
	readable.WriteString(`{"class":"DEVICES","devices":[{"class":"DEVICE","path":"/dev/gps0","driver":"SiRF","subtype":"9\u0006GSD4e_4.1.2-P1_RPATCH.05-F-GPS-4R-1510281 11/03/2015 307","activated":"2020-03-03T21:54:51.357Z","flags":1,"native":1,"bps":4800,"parity":"N","stopbits":1,"cycle":1.00}]}`)
	readable.WriteString(`{"class":"WATCH","enable":true,"json":true,"nmea":false,"raw":0,"scaled":false,"timing":false,"split24":false,"pps":false}`)
	readable.WriteString(`{"class":"TPV","device":"/dev/gps0","status":2,"mode":3,"time":"2020-03-03T21:54:49.000Z","ept":0.005,"lat":40.7828687,"lon":-73.9675438,"alt":166.395,"epx":3.271,"epy":4.863,"epv":14.780,"track":0.0000,"speed":0.000,"climb":0.000,"eps":0.19,"epc":0.58}`)

	writeable := bytes.Buffer{}

	conn := &MockConn{
		readableData: &readable,
		writtenData:  &writeable,
	}

	g := gpsd.New(conn)

	g.Watch()

	assert.Equal(t, `?WATCH={"enable":true,"json":true}`, writeable.String())

	messages, errors := g.Scan(gpsd.MessageTypeTimePositionVelocity | gpsd.MessageTypeVersion)

	var tpvCalled bool
	var versionCalled bool

	h := &gpsd.Handler{}
	h = h.WithTimePositionVelocity(func(rpt gpsd.TimePositionVelocityReport) {
		assert.Equal(t, rpt, gpsd.TimePositionVelocityReport{
			Class:  "TPV",
			Tag:    "",
			Device: "/dev/gps0",
			Mode:   gpsd.Mode3D,
			Time:   time.Date(2020, 03, 03, 21, 54, 49, 0, time.UTC),
			Ept:    0.005,
			Lat:    40.7828687,
			Lon:    -73.9675438,
			Alt:    166.395,
			Epx:    3.271,
			Epy:    4.863,
			Epv:    14.78,
			Track:  0,
			Speed:  0,
			Climb:  0,
			Epd:    0,
			Eps:    0.19,
			Epc:    0.58,
		})
		tpvCalled = true
	})
	h = h.WithVersion(func(rpt gpsd.VersionReport) {
		assert.Equal(t, rpt, gpsd.VersionReport{
			Class:      "VERSION",
			Release:    "3.17",
			Rev:        "3.17",
			ProtoMajor: 3,
			ProtoMinor: 12,
			Remote:     "",
		})
		versionCalled = true
	})

	gpsd.ScanLoop(messages, errors, h, func(err error) {
		t.Log(err.Error())
		t.FailNow()
	})

	assert.Equal(t, true, tpvCalled)
	assert.Equal(t, true, versionCalled)
}

func ExampleGPSD() {
	// Typically you would use net.Dial("tcp", "localhost:2947") instead of making a fake conn here.
	readable := bytes.Buffer{}

	readable.WriteString(`{"class":"VERSION","release":"3.17","rev":"3.17","proto_major":3,"proto_minor":12}`)
	readable.WriteString(`{"class":"TPV","device":"/dev/gps0","status":2,"mode":3,"time":"2020-03-03T21:54:49.000Z","ept":0.005,"lat":40.7828687,"lon":-73.9675438,"alt":166.395,"epx":3.271,"epy":4.863,"epv":14.780,"track":0.0000,"speed":0.000,"climb":0.000,"eps":0.19,"epc":0.58}`)

	writeable := bytes.Buffer{}

	conn := &MockConn{
		readableData: &readable,
		writtenData:  &writeable,
	}

	g := gpsd.New(conn)
	g.Watch()

	time.Sleep(time.Millisecond)

	messages, errors := g.Scan(gpsd.MessageTypeTimePositionVelocity | gpsd.MessageTypeVersion)

	h := &gpsd.Handler{}
	h = h.WithTimePositionVelocity(func(rpt gpsd.TimePositionVelocityReport) {
		fmt.Println(fmt.Sprintf("Location: %0.7f, %0.7f", rpt.Lat, rpt.Lon))

		// only care about the first location event, so stop reading from gpsd when we get it
		g.Stop()
	}).WithVersion(func(rpt gpsd.VersionReport) {
		fmt.Println(fmt.Sprintf("Connected to GPSD version %s", rpt.Release))
	})

	gpsd.ScanLoop(messages, errors, h, func(err error) {
		panic(err.Error())
	})

	// Output:
	// Connected to GPSD version 3.17
	// Location: 40.7828687, -73.9675438
}
