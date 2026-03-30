package gomavlinkdroneapi

import (
	"github.com/mavlink/MAVSDK-Go/Sources/action"
	"github.com/mavlink/MAVSDK-Go/Sources/core"
	"github.com/mavlink/MAVSDK-Go/Sources/geofence"
)

type API interface {
	Connect()
	Arm() (*action.ArmResponse, error)
	Takeoff() (*action.TakeoffResponse, error)
	Land() (*action.LandResponse, error)
	ConnectionState() (<-chan *core.ConnectionState, error)
	UploadGeofence(geofenceData *geofence.GeofenceData) (*geofence.UploadGeofenceResponse, error)
}

// go mod init github.com/GolangToolKits/GoMavlinkDroneAPI
