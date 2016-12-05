package event

// event bus events handled by the device package
const (
	DeviceRegister   Type = "device:register"
	DeviceAddToZone  Type = "device:toZone"
	DeviceBind       Type = "device:bind"
	DeviceUnregister Type = "device:unregister"
)
