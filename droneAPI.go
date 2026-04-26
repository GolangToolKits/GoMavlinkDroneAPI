package gomavlinkdroneapi

import (
	"github.com/mavlink/MAVSDK-Go/Sources/action"
	"github.com/mavlink/MAVSDK-Go/Sources/core"
	"github.com/mavlink/MAVSDK-Go/Sources/geofence"
)

type DroneAPI struct {
	// Port         string
	// MavsdkServer string
	Baud        int
	Device      string
	OutSystemID int
	Drone       *Drone
}

func (s *DroneAPI) Connect() error {
	// drone := &Drone{port: "50051", mavsdkServer: "127.0.0.1"}
	s.Drone = &Drone{baud: s.Baud, device: s.Device, outSystemID: s.OutSystemID}
	err := s.Drone.Connect()
	return err
}

func (s *DroneAPI) Arm() (*action.ArmResponse, error) {
	// ar, err := s.Drone.action.Arm(context.Background())
	// return ar, err
	return nil, nil
}

func (s *DroneAPI) Takeoff() (*action.TakeoffResponse, error) {
	// tr, err := s.Drone.action.Takeoff(context.Background())
	// return tr, err
	return nil, nil
}

func (s *DroneAPI) Land() (*action.LandResponse, error) {
	// lr, err := s.Drone.action.Land(context.Background())
	// return lr, err
	return nil, nil
}

func (s *DroneAPI) ConnectionState() (<-chan *core.ConnectionState, error) {
	// cs, err := s.Drone.core.ConnectionState(context.Background())
	// return cs, err
	return nil, nil
}

func (s *DroneAPI) UploadGeofence(geofenceData *geofence.GeofenceData) (*geofence.UploadGeofenceResponse, error) {
	// gfr, err := s.Drone.geofence.UploadGeofence(context.Background(), geofenceData)
	// return gfr, err
	return nil, nil
}

func (s *DroneAPI) Close() {
	s.Drone.Close()
}

func (s *DroneAPI) New() API {
	return s
}

// lat := 47.3977508
// lon := 8.5456074
// p1 := &geofence.Point{
// 	LatitudeDeg:  lat - 0.0001,
// 	LongitudeDeg: lon - 0.0001,
// }
// p2 := &geofence.Point{
// 	LatitudeDeg:  lat + 0.0001,
// 	LongitudeDeg: lon - 0.0001,
// }
// p3 := &geofence.Point{
// 	LatitudeDeg:  lat + 0.0001,
// 	LongitudeDeg: lon + 0.0001,
// }
// p4 := &geofence.Point{
// 	LatitudeDeg:  lat - 0.0001,
// 	LongitudeDeg: lon + 0.0001,
// }
// // this is not a test or verification package. this only checks the sanity of geofence api
// polygon := &geofence.Polygon{
// 	Points:    []*geofence.Point{p1, p2, p3, p4},
// 	FenceType: geofence.FenceType_FENCE_TYPE_EXCLUSION}
// response, err := drone.geofence.UploadGeofence(context.Background(), &geofence.GeofenceData{
// 	Polygons: []*geofence.Polygon{polygon},
// })
// if err != nil {
// 	log.Print(err.Error())
// 	os.Exit(1)
// }
// log.Printf("response %v", response)
