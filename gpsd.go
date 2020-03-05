package gpsd

import (
	"encoding/json"
	"io"
)

// MessageType is a flag used to tell GPSD.Scan what message types to scan for.
type MessageType uint16

const (
	// MessageTypeTimePositionVelocity represents the "TimePositionVelocity" message.
	MessageTypeTimePositionVelocity MessageType = 1 << iota
	// MessageTypeVersion represents the "Version" message.
	MessageTypeVersion
	// MessageTypeSkyView represents the "SkyView" message.
	MessageTypeSkyView
	// MessageTypePseudorangeNoise represents the "PseudorangeNoise" message.
	MessageTypePseudorangeNoise
	// MessageTypeVehicleAttitude represents the "VehicleAttitude" message.
	MessageTypeVehicleAttitude
	// MessageTypeDevices represents the "Devices" message.
	MessageTypeDevices
	// MessageTypePulsePerSecond represents the "PulsePerSecond" message.
	MessageTypePulsePerSecond
	// MessageTypeError represents the "Error" message.
	MessageTypeError
)

// Messages wraps all the channels used to receive messages from gpsd.
type Messages struct {
	tpv chan TimePositionVelocityReport
	ver chan VersionReport
	sky chan SkyViewReport
	gst chan PseudorangeNoiseReport
	att chan VehicleAttitudeReport
	dev chan DevicesReport
	pps chan PulsePerSecondReport
	err chan ErrorReport
}

func (m *Messages) init(mt MessageType) {
	if mt&MessageTypeTimePositionVelocity != 0 {
		m.tpv = make(chan TimePositionVelocityReport, 100)
	}
	if mt&MessageTypeVersion != 0 {
		m.ver = make(chan VersionReport, 100)
	}
	if mt&MessageTypeSkyView != 0 {
		m.sky = make(chan SkyViewReport, 100)
	}
	if mt&MessageTypePseudorangeNoise != 0 {
		m.gst = make(chan PseudorangeNoiseReport, 100)
	}
	if mt&MessageTypeVehicleAttitude != 0 {
		m.att = make(chan VehicleAttitudeReport, 100)
	}
	if mt&MessageTypeDevices != 0 {
		m.dev = make(chan DevicesReport, 100)
	}
	if mt&MessageTypePulsePerSecond != 0 {
		m.pps = make(chan PulsePerSecondReport, 100)
	}
	if mt&MessageTypeError != 0 {
		m.err = make(chan ErrorReport, 100)
	}
}

func (m *Messages) close() {
	if m.tpv != nil {
		close(m.tpv)
	}

	if m.ver != nil {
		close(m.ver)
	}

	if m.sky != nil {
		close(m.sky)
	}

	if m.gst != nil {
		close(m.gst)
	}

	if m.att != nil {
		close(m.att)
	}

	if m.dev != nil {
		close(m.dev)
	}

	if m.pps != nil {
		close(m.pps)
	}

	if m.err != nil {
		close(m.err)
	}
}

// TimePositionVelocity returns the channel that receives TimePositionVelocity reports.
func (m *Messages) TimePositionVelocity() <-chan TimePositionVelocityReport {
	return m.tpv
}

// Version returns the channel that receives Version reports.
func (m *Messages) Version() <-chan VersionReport {
	return m.ver
}

// SkyView returns the channel that receives SkyView reports.
func (m *Messages) SkyView() <-chan SkyViewReport {
	return m.sky
}

// PseudorangeNoise returns the channel that receives PseudorangeNoise reports.
func (m *Messages) PseudorangeNoise() <-chan PseudorangeNoiseReport {
	return m.gst
}

// VehicleAttitude returns the channel that receives VehicleAttitude reports.
func (m *Messages) VehicleAttitude() <-chan VehicleAttitudeReport {
	return m.att
}

// Devices returns the channel that receives Devices reports.
func (m *Messages) Devices() <-chan DevicesReport {
	return m.dev
}

// PulsePerSecond returns the channel that receives PulsePerSecond reports.
func (m *Messages) PulsePerSecond() <-chan PulsePerSecondReport {
	return m.pps
}

// Error returns the channel that receives Error reports.
func (m *Messages) Error() <-chan ErrorReport {
	return m.err
}

type message struct {
	Class string `json:"class"`
}

// GPSD represents the logic to receive messages from gpsd.
type GPSD struct {
	conn io.ReadWriter
	stop bool
}

// New creates a new GPSD with the given connection. You can use net.Dial to
// connect to gpsd and pass the *net.Conn in here.
func New(conn io.ReadWriter) *GPSD {
	return &GPSD{
		conn: conn,
	}
}

// Watch will send the WATCH command to gpsd.
func (g *GPSD) Watch() error {
	_, err := g.conn.Write([]byte(`?WATCH={"enable":true,"json":true}`))
	return err
}

// Scan will start scanning the connection to gpsd for the message types
// specified.
func (g *GPSD) Scan(mt MessageType) (*Messages, <-chan error) {
	messages := &Messages{}
	messages.init(mt)

	errChan := make(chan error, 100)

	go func(m *Messages, errChan chan error) {
		g.watchLoop(m, errChan)

		messages.close()
		close(errChan)
	}(messages, errChan)

	return messages, errChan
}

// Stop will stop scanning for new messages on the connection to gpsd.
func (g *GPSD) Stop() {
	g.stop = true
}

func (g *GPSD) watchLoop(m *Messages, errChan chan error) {
	dec := json.NewDecoder(g.conn)

	for !g.stop {
		var raw json.RawMessage
		err := dec.Decode(&raw)
		if err != nil {
			if err != io.EOF {
				errChan <- err
			}

			g.stop = true
			continue
		}

		var msg message
		err = json.Unmarshal(raw, &msg)
		if err != nil {
			errChan <- err
			continue
		}

		switch msg.Class {
		case "TPV":
			if m.tpv != nil {
				var rpt TimePositionVelocityReport
				err = json.Unmarshal(raw, &rpt)
				if err != nil {
					errChan <- err
					continue
				}

				m.tpv <- rpt
			}
		case "SKY":
			if m.sky != nil {
				var rpt SkyViewReport
				err = json.Unmarshal(raw, &rpt)
				if err != nil {
					errChan <- err
					continue
				}

				m.sky <- rpt
			}
		case "GST":
			if m.gst != nil {
				var rpt PseudorangeNoiseReport
				err = json.Unmarshal(raw, &rpt)
				if err != nil {
					errChan <- err
					continue
				}

				m.gst <- rpt
			}
		case "ATT":
			if m.att != nil {
				var rpt VehicleAttitudeReport
				err = json.Unmarshal(raw, &rpt)
				if err != nil {
					errChan <- err
					continue
				}

				m.att <- rpt
			}
		case "VERSION":
			if m.ver != nil {
				var rpt VersionReport
				err = json.Unmarshal(raw, &rpt)
				if err != nil {
					errChan <- err
					continue
				}

				m.ver <- rpt
			}
		case "DEVICES":
			if m.dev != nil {
				var rpt DevicesReport
				err = json.Unmarshal(raw, &rpt)
				if err != nil {
					errChan <- err
					continue
				}

				m.dev <- rpt
			}
		case "PPS":
			if m.pps != nil {
				var rpt PulsePerSecondReport
				err = json.Unmarshal(raw, &rpt)
				if err != nil {
					errChan <- err
					continue
				}

				m.pps <- rpt
			}
		case "ERROR":
			if m.err != nil {
				var rpt ErrorReport
				err = json.Unmarshal(raw, &rpt)
				if err != nil {
					errChan <- err
					continue
				}

				m.err <- rpt
			}
		}
	}
}

// ScanLoop will monitor all the channels on the Messages wrapper given, calling
// the functions on the given ReportHandler or ErrorHandler for each message
// received.
func ScanLoop(m *Messages, errs <-chan error, h ReportHandler, eh ErrorHandler) {
	tpvRpts := m.TimePositionVelocity()
	verRpts := m.Version()
	skyRpts := m.SkyView()
	gstRpts := m.PseudorangeNoise()
	attRpts := m.VehicleAttitude()
	devRpts := m.Devices()
	ppsRpts := m.PulsePerSecond()
	errRpts := m.Error()

	for {
		select {
		case err, ok := <-errs:
			if !ok {
				errs = nil
				break
			}

			eh(err)

		case rpt, ok := <-tpvRpts:
			if !ok {
				tpvRpts = nil
				break
			}

			h.TimePositionVelocity(rpt)

		case rpt, ok := <-verRpts:
			if !ok {
				verRpts = nil
				break
			}

			h.Version(rpt)

		case rpt, ok := <-skyRpts:
			if !ok {
				skyRpts = nil
				break
			}

			h.SkyView(rpt)

		case rpt, ok := <-gstRpts:
			if !ok {
				gstRpts = nil
				break
			}

			h.PseudorangeNoise(rpt)

		case rpt, ok := <-attRpts:
			if !ok {
				attRpts = nil
				break
			}

			h.VehicleAttitude(rpt)

		case rpt, ok := <-devRpts:
			if !ok {
				devRpts = nil
				break
			}

			h.Devices(rpt)

		case rpt, ok := <-ppsRpts:
			if !ok {
				ppsRpts = nil
				break
			}

			h.PulsePerSecond(rpt)

		case rpt, ok := <-errRpts:
			if !ok {
				errRpts = nil
				break
			}

			h.Error(rpt)
		}

		if errs == nil && tpvRpts == nil && verRpts == nil &&
			skyRpts == nil && gstRpts == nil && attRpts == nil &&
			devRpts == nil && ppsRpts == nil && errRpts == nil {
			// all channels closed
			break
		}
	}
}
