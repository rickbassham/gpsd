package gpsd

import "time"

// Mode is the NMEA mode of a TimePositionVelocity report.
type Mode byte

const (
	// NoValueSeen indicates no data has been received yet.
	NoValueSeen Mode = 0
	// NoFix indicates fix has not been required yet.
	NoFix Mode = 1
	// Mode2D represents quality of the fix.
	Mode2D Mode = 2
	// Mode3D represents quality of the fix.
	Mode3D Mode = 3
)

// TimePositionVelocityReport is a TPV report.
type TimePositionVelocityReport struct {
	Class  string    `json:"class"`
	Tag    string    `json:"tag"`
	Device string    `json:"device"`
	Mode   Mode      `json:"mode"`
	Time   time.Time `json:"time"`
	Ept    float64   `json:"ept"`
	Lat    float64   `json:"lat"`
	Lon    float64   `json:"lon"`
	Alt    float64   `json:"alt"`
	Epx    float64   `json:"epx"`
	Epy    float64   `json:"epy"`
	Epv    float64   `json:"epv"`
	Track  float64   `json:"track"`
	Speed  float64   `json:"speed"`
	Climb  float64   `json:"climb"`
	Epd    float64   `json:"epd"`
	Eps    float64   `json:"eps"`
	Epc    float64   `json:"epc"`
}

// SkyViewReport is a SKY report.
type SkyViewReport struct {
	Class      string      `json:"class"`
	Tag        string      `json:"tag"`
	Device     string      `json:"device"`
	Time       time.Time   `json:"time"`
	Xdop       float64     `json:"xdop"`
	Ydop       float64     `json:"ydop"`
	Vdop       float64     `json:"vdop"`
	Tdop       float64     `json:"tdop"`
	Hdop       float64     `json:"hdop"`
	Pdop       float64     `json:"pdop"`
	Gdop       float64     `json:"gdop"`
	Satellites []Satellite `json:"satellites"`
}

// PseudorangeNoiseReport is a GST report.
type PseudorangeNoiseReport struct {
	Class  string    `json:"class"`
	Tag    string    `json:"tag"`
	Device string    `json:"device"`
	Time   time.Time `json:"time"`
	Rms    float64   `json:"rms"`
	Major  float64   `json:"major"`
	Minor  float64   `json:"minor"`
	Orient float64   `json:"orient"`
	Lat    float64   `json:"lat"`
	Lon    float64   `json:"lon"`
	Alt    float64   `json:"alt"`
}

// VehicleAttitudeReport is an ATT report.
type VehicleAttitudeReport struct {
	Class       string    `json:"class"`
	Tag         string    `json:"tag"`
	Device      string    `json:"device"`
	Time        time.Time `json:"time"`
	Heading     float64   `json:"heading"`
	MagSt       string    `json:"mag_st"`
	Pitch       float64   `json:"pitch"`
	PitchSt     string    `json:"pitch_st"`
	Yaw         float64   `json:"yaw"`
	YawSt       string    `json:"yaw_st"`
	Roll        float64   `json:"roll"`
	RollSt      string    `json:"roll_st"`
	Dip         float64   `json:"dip"`
	MagLen      float64   `json:"mag_len"`
	MagX        float64   `json:"mag_x"`
	MagY        float64   `json:"mag_y"`
	MagZ        float64   `json:"mag_z"`
	AccLen      float64   `json:"acc_len"`
	AccX        float64   `json:"acc_x"`
	AccY        float64   `json:"acc_y"`
	AccZ        float64   `json:"acc_z"`
	GyroX       float64   `json:"gyro_x"`
	GyroY       float64   `json:"gyro_y"`
	Depth       float64   `json:"depth"`
	Temperature float64   `json:"temperature"`
}

// VersionReport is a VERSION report.
type VersionReport struct {
	Class      string `json:"class"`
	Release    string `json:"release"`
	Rev        string `json:"rev"`
	ProtoMajor int    `json:"proto_major"`
	ProtoMinor int    `json:"proto_minor"`
	Remote     string `json:"remote"`
}

// DevicesReport is a DEVICES report.
type DevicesReport struct {
	Class   string   `json:"class"`
	Devices []Device `json:"devices"`
	Remote  string   `json:"remote"`
}

// Device is a single device connected to gpsd.
type Device struct {
	Class     string  `json:"class"`
	Path      string  `json:"path"`
	Activated string  `json:"activated"`
	Flags     int     `json:"flags"`
	Driver    string  `json:"driver"`
	Subtype   string  `json:"subtype"`
	Bps       int     `json:"bps"`
	Parity    string  `json:"parity"`
	Stopbits  int     `json:"stopbits"`
	Native    int     `json:"native"`
	Cycle     float64 `json:"cycle"`
	Mincycle  float64 `json:"mincycle"`
}

// PulsePerSecondReport is a PPS report.
type PulsePerSecondReport struct {
	Class      string  `json:"class"`
	Device     string  `json:"device"`
	RealSec    float64 `json:"real_sec"`
	RealMusec  float64 `json:"real_musec"`
	ClockSec   float64 `json:"clock_sec"`
	ClockMusec float64 `json:"clock_musec"`
}

// ErrorReport is an ERROR report.
type ErrorReport struct {
	Class   string `json:"class"`
	Message string `json:"message"`
}

// Satellite represents a satellite in the SKY report.
type Satellite struct {
	PRN  float64 `json:"PRN"`
	Az   float64 `json:"az"`
	El   float64 `json:"el"`
	Ss   float64 `json:"ss"`
	Used bool    `json:"used"`
}
