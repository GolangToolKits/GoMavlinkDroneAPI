package gomavlinkdroneapi

import (
	"fmt"

	"github.com/mavlink/MAVSDK-Go/Sources/action"
	"github.com/mavlink/MAVSDK-Go/Sources/camera"
	"github.com/mavlink/MAVSDK-Go/Sources/core"
	"github.com/mavlink/MAVSDK-Go/Sources/geofence"
	"github.com/mavlink/MAVSDK-Go/Sources/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Drone struct {
	port         string
	mavsdkServer string
	action       action.ServiceImpl
	core         core.ServiceImpl
	telemetry    telemetry.ServiceImpl
	geofence     geofence.ServiceImpl
	camera       camera.ServiceImpl
}

// Connect Starts a mavsdk server and create a connection to it
func (s *Drone) Connect() {
	grpcConnection := s.connectToMAVSDKServer()
	s.InitPlugins(grpcConnection)

}

// InitPlugins initializes all the plugins
func (s *Drone) InitPlugins(cc *grpc.ClientConn) {

	s.telemetry = telemetry.ServiceImpl{
		Client: telemetry.NewTelemetryServiceClient(cc),
	}
	s.core = core.ServiceImpl{
		Client: core.NewCoreServiceClient(cc),
	}
	s.action = action.ServiceImpl{
		Client: action.NewActionServiceClient(cc),
	}
	// s.action = action.ServiceImpl{
	// 	Client: action.NewActionServiceClient(cc),
	// }
	s.geofence = geofence.ServiceImpl{
		Client: geofence.NewGeofenceServiceClient(cc),
	}
	s.camera = camera.ServiceImpl{
		Client: camera.NewCameraServiceClient(cc),
	}
}

func (s *Drone) connectToMAVSDKServer() *grpc.ClientConn {
	dialoption := grpc.WithTransportCredentials(insecure.NewCredentials())

	serverAddr := s.mavsdkServer + ":" + s.port
	cc, err := grpc.NewClient(serverAddr, dialoption)
	if err != nil {
		fmt.Printf("Error while dialing %v", err)
	}
	grpc.ConnectionTimeout(5)
	return cc
}
