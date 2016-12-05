package device

import "net"

//Type represents the type of connected device
type Type string

// Predefined device types
const (
	Audio   Type = "audio"
	Display Type = "display"
)

// Device holds all information about conected devices
type Device struct {
	ID             string   `json:"id"`
	IP             string   `json:"ip"`
	Port           string   `json:"port"`
	Live           bool     `json:"live"`
	Type           Type     `json:"type"`
	Bind           bool     `json:"bind"`
	BindTimeout    int      `json:"bindTimeout"`
	Zones          []string `json:"zones"`
	tcpConn        net.Conn
	tpcIsConnected bool
}
