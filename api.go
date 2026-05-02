package gomavlinkdroneapi

import (
	"context"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

type API interface {

	//connect to the vehicle
	ConnectSerial(serialDevice string, baud int, outSystemID byte) error
	ConnectUDPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectUDPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectUDPBroadcast(broadcastAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectTCPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectTCPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectCustomClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error
	ConnectCustomServer(listenAddress string, outVersion gomavlib.Version, outSystemID byte) error

	// IsDroneConnected() (bool, error)
	IsDroneConnected(ctx context.Context) (bool, error)

	Close()

	//Get Connected Vehicle ids
	// GetConnectedVehicle() (targetSystem byte, targetComponent byte)
	GetConnectedVehicle(ctx context.Context) (targetSystem byte, targetComponent byte)

	// listen to the vehicle
	// ListenToDroneEvents() chan gomavlib.Event

	//upload missions to the vehicle
	// UploadMission(droneChannel *gomavlib.Channel, mission []ardupilotmega.MessageMissionItemInt) bool
	// UploadMission(missionItems *[]ardupilotmega.MessageMissionItemInt, targetSystem uint8, targetComponent uint8) bool
	// DownloadMissions(ctx context.Context, targetSystem uint8, targetComponent uint8) (map[uint16]*ardupilotmega.MessageMissionItemInt, error)
	//DownloadMissions(targetSystem uint8, targetComponent uint8) (map[uint16]*ardupilotmega.MessageMissionItemInt, error)
	UploadMission(ctx context.Context, missionItems []ardupilotmega.MessageMissionItemInt, targetSystem uint8, targetComponent uint8) error
	DownloadMissions(ctx context.Context, targetSystem uint8, targetComponent uint8) ([]*ardupilotmega.MessageMissionItemInt, error)
	// OverrideMissionAndHover(command *common.MessageCommandLong) error
	OverrideMissionAndHover(ctx context.Context, command *common.MessageCommandLong) error
	// GetDroneChannel() *gomavlib.Channel
	GetDroneChannel(ctx context.Context) *gomavlib.Channel

	// flying the vehicle
	//ArmDisarm(command float32, targetSystem uint8, targetComponent uint8) error
	ArmDisarm(ctx context.Context, arm bool, targetSystem uint8, targetComponent uint8) error
	// Takeoff(takeOffCommand *common.MessageCommandLong) error
	Takeoff(ctx context.Context, altitude float32, targetSystem uint8, targetComponent uint8) error
	// Land(landCommand *common.MessageCommandLong) error
	Land(ctx context.Context, targetSystem uint8, targetComponent uint8) error
	// Move(moveMessage *common.MessageSetPositionTargetLocalNed) error
	Move(targetSystem uint8, targetComponent uint8, moveMessage *common.MessageSetPositionTargetLocalNed) error
	ReturnHome(ctx context.Context, targetSystem uint8, targetComponent uint8) error
	// ReturnHome(rtlCommand *common.MessageCommandLong) error
	// AcknowledgeCommand(commandToCheck common.MAV_CMD) bool

	// ClearGeofence(targetSystem uint8, targetComponent uint8) bool
	ClearGeofence(ctx context.Context, targetSystem uint8, targetComponent uint8) bool
	UploadGeofence(ctx context.Context, newFence []GeoFence, targetSystem uint8, targetComponent uint8) (bool, error)
	// UploadGeofence(newFence *[]GeoFence, targetSystem uint8, targetComponent uint8) (bool, error)
}

// go mod init github.com/GolangToolKits/GoMavlinkDroneAPI
