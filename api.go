package gomavlinkdroneapi

import (
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
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
	IsDroneConnected() (bool, error)
	// UploadMission(droneChannel *gomavlib.Channel, mission []ardupilotmega.MessageMissionItemInt) bool
	UploadMission(missionItems []ardupilotmega.MessageMissionItemInt, targetSystem uint8, targetComponent uint8) bool
	// GetDroneChannel() *gomavlib.Channel
	Close()
	ArmDisarm(command float32, targetSystem uint8, targetComponent uint8) error
	Takeoff(takeOffCommand *common.MessageCommandLong) error
	Land(landCommand *common.MessageCommandLong) error
	AcknowledgeCommand(commandToCheck common.MAV_CMD) bool
	// UploadGeofence(geofenceData *geofence.GeofenceData) (*geofence.UploadGeofenceResponse, error)
}

// go mod init github.com/GolangToolKits/GoMavlinkDroneAPI
