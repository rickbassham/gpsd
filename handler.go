package gpsd

// Handler is a helper struct that implements ReportHandler.
type Handler struct {
	tpv func(TimePositionVelocityReport)
	ver func(VersionReport)
	sky func(SkyViewReport)
	gst func(PseudorangeNoiseReport)
	att func(VehicleAttitudeReport)
	dev func(DevicesReport)
	pps func(PulsePerSecondReport)
	err func(ErrorReport)
}

// WithTimePositionVelocity adds a function handler for TimePositionVelocityReport messages.
func (h *Handler) WithTimePositionVelocity(fn func(TimePositionVelocityReport)) *Handler {
	h.tpv = fn
	return h
}

// WithVersion adds a function handler for VersionReport messages.
func (h *Handler) WithVersion(fn func(VersionReport)) *Handler {
	h.ver = fn
	return h
}

// WithSkyView adds a function handler for SkyViewReport messages.
func (h *Handler) WithSkyView(fn func(SkyViewReport)) *Handler {
	h.sky = fn
	return h
}

// WithPseudorangeNoise adds a function handler for PseudorangeNoiseReport messages.
func (h *Handler) WithPseudorangeNoise(fn func(PseudorangeNoiseReport)) *Handler {
	h.gst = fn
	return h
}

// WithVehicleAttitude adds a function handler for VehicleAttitudeReport messages.
func (h *Handler) WithVehicleAttitude(fn func(VehicleAttitudeReport)) *Handler {
	h.att = fn
	return h
}

// WithDevices adds a function handler for DevicesReport messages.
func (h *Handler) WithDevices(fn func(DevicesReport)) *Handler {
	h.dev = fn
	return h
}

// WithPulsePerSecond adds a function handler for PulsePerSecondReport messages.
func (h *Handler) WithPulsePerSecond(fn func(PulsePerSecondReport)) *Handler {
	h.pps = fn
	return h
}

// WithError adds a function handler for ErrorReport messages.
func (h *Handler) WithError(fn func(ErrorReport)) *Handler {
	h.err = fn
	return h
}

// TimePositionVelocity calls the given handler.
func (h *Handler) TimePositionVelocity(rpt TimePositionVelocityReport) {
	h.tpv(rpt)
}

// Version calls the given handler.
func (h *Handler) Version(rpt VersionReport) {
	h.ver(rpt)
}

// SkyView calls the given handler.
func (h *Handler) SkyView(rpt SkyViewReport) {
	h.sky(rpt)
}

// PseudorangeNoise calls the given handler.
func (h *Handler) PseudorangeNoise(rpt PseudorangeNoiseReport) {
	h.gst(rpt)
}

// VehicleAttitude calls the given handler.
func (h *Handler) VehicleAttitude(rpt VehicleAttitudeReport) {
	h.att(rpt)
}

// Devices calls the given handler.
func (h *Handler) Devices(rpt DevicesReport) {
	h.dev(rpt)
}

// PulsePerSecond calls the given handler.
func (h *Handler) PulsePerSecond(rpt PulsePerSecondReport) {
	h.pps(rpt)
}

// Error calls the given handler.
func (h *Handler) Error(rpt ErrorReport) {
	h.err(rpt)
}

// ReportHandler is an interface that can be implemented to use with ScanLoop.
type ReportHandler interface {
	TimePositionVelocity(TimePositionVelocityReport)
	Version(VersionReport)
	SkyView(SkyViewReport)
	PseudorangeNoise(PseudorangeNoiseReport)
	VehicleAttitude(VehicleAttitudeReport)
	Devices(DevicesReport)
	PulsePerSecond(PulsePerSecondReport)
	Error(ErrorReport)
}

// ErrorHandler defines the function signature to handle errors.
type ErrorHandler func(error)
