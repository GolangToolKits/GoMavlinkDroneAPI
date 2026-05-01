package gomavlinkdroneapi

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

type GeoFence struct {
	Lat int32
	Lon int32
}

func (s *DroneAPI) ClearGeofence(targetSystem uint8, targetComponent uint8) bool {
	clearMsg := &common.MessageMissionClearAll{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		MissionType:     common.MAV_MISSION_TYPE_FENCE, // Ensures only fences are wiped
	}

	time.Sleep(1 * time.Second) // Wait for connection
	s.drone.node.WriteMessageAll(clearMsg)
	log.Println("Request to clear old geofences sent.")

	// Use a clean range loop to process incoming events
	for evt := range s.drone.node.Events() {

		// Validate incoming MAVLink frame
		if e, ok := evt.(*gomavlib.EventFrame); ok {
			if msg, ok := e.Message().(*common.MessageMissionAck); ok {
				if msg.MissionType == common.MAV_MISSION_TYPE_FENCE {
					if msg.Type == common.MAV_MISSION_ACCEPTED {
						log.Println("Geofence cleared successfully!")
						return true
					}
					log.Printf("Failed to clear. Error code: %d\n", msg.Type)
					return false
				}
			}
		}
	}

	// This is only hit if gomavlib stops and closes the Events() channel
	log.Println("Node events channel closed.")
	return false
}

// send new geofense data
// Example
//
//	newFence: &[]gomavlinkdroneapi.GeoFence{
//					{Lat: 473977418, Lon: 85455939},
//					{Lat: 473977418, Lon: 85465939},
//					{Lat: 473987418, Lon: 85465939},
//					{Lat: 473987418, Lon: 85455939},
//				},
func (s *DroneAPI) UploadGeofence(newFence *[]GeoFence, targetSystem uint8, targetComponent uint8) (bool, error) {
	log.Println("Commencing new fence stream...")
	time.Sleep(1 * time.Second)

	// Declare boundary volume size
	s.drone.node.WriteMessageAll(&common.MessageMissionCount{
		TargetSystem: targetSystem, TargetComponent: targetComponent,
		Count: uint16(len(*newFence)), MissionType: common.MAV_MISSION_TYPE_FENCE,
	})

	// Use a clean range loop to process incoming events indefinitely
	for evt := range s.drone.node.Events() {
		e, ok := evt.(*gomavlib.EventFrame)
		if !ok {
			continue
		}

		// Autopilot requests a point
		if msg, ok := e.Message().(*common.MessageMissionRequestInt); ok {
			if msg.MissionType == common.MAV_MISSION_TYPE_FENCE {
				seq := msg.Seq

				if int(seq) < len(*newFence) {
					log.Printf("Uploading requested point #%d\n", seq)
					s.drone.node.WriteMessageAll(&common.MessageMissionItemInt{
						TargetSystem:    targetSystem,
						TargetComponent: targetComponent,
						Seq:             seq,
						Frame:           common.MAV_FRAME_GLOBAL,
						Command:         common.MAV_CMD_NAV_FENCE_POLYGON_VERTEX_INCLUSION,
						Param1:          float32(len(*newFence)),
						X:               (*newFence)[seq].Lat,
						Y:               (*newFence)[seq].Lon,
						MissionType:     common.MAV_MISSION_TYPE_FENCE,
					})
				}
			}
		}

		// Autopilot acknowledges completion
		if msg, ok := e.Message().(*common.MessageMissionAck); ok {
			if msg.MissionType == common.MAV_MISSION_TYPE_FENCE {
				if msg.Type == common.MAV_MISSION_ACCEPTED {
					log.Println("New dynamic geofence live and locked!")
					return true, nil
				}
				log.Printf("Rejected with code %d\n", msg.Type)
				// Fixed the rune bug here: replaced string(rune(msg.Type)) with fmt.Sprintf
				return false, fmt.Errorf("rejected with code %d", msg.Type)
			}
		}
	}

	// This is only hit if gomavlib stops and closes the Events() channel
	log.Println("Node events channel closed unexpectedly.")
	return false, errors.New("node events channel closed unexpectedly")
}
