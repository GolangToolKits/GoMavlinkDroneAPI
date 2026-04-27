package gomavlinkdroneapi

import (
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

type API interface {
	ConnectSerial(serialDevice string, baud int, outSystemID byte) error
	ConnectUDPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectUDPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectUDPBroadcast(broadcastAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectTCPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectTCPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectCustomClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectCustomServer(listenAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ListenToDroneEvents() chan gomavlib.Event
	IsDroneConnected() bool
	UploadMission(droneChannel *gomavlib.Channel, mission []ardupilotmega.MessageMissionItemInt) bool
	GetDroneChannel() *gomavlib.Channel
	Close()
	ArmDisarm(command float32, targetSystem uint8, targetComponent uint8) error
	// Takeoff() (*action.TakeoffResponse, error)
	// Land() (*action.LandResponse, error)
	// ConnectionState() (<-chan *core.ConnectionState, error)
	// UploadGeofence(geofenceData *geofence.GeofenceData) (*geofence.UploadGeofenceResponse, error)
}

// go mod init github.com/GolangToolKits/GoMavlinkDroneAPI
