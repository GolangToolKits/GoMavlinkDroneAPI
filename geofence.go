package gomavlinkdroneapi

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

type GeoFence struct {
	Lat int32
	Lon int32
}

func (s *DroneAPI) ClearGeofence(ctx context.Context, targetSystem uint8, targetComponent uint8) bool {
	clearMsg := &common.MessageMissionClearAll{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		MissionType:     common.MAV_MISSION_TYPE_FENCE,
	}

	// Use WriteMessageAll if you haven't captured a specific channel yet,
	// otherwise, use s.drone.node.WriteMessageTo(channel, clearMsg)
	s.drone.node.WriteMessageAll(clearMsg)
	log.Printf("Request to clear geofence sent to System %d, Component %d", targetSystem, targetComponent)

	for {
		select {
		case <-ctx.Done():
			log.Println("ClearGeofence operation timed out.")
			return false
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return false
			}

			if frm, ok := evt.(*gomavlib.EventFrame); ok {
				// 1. Verify the message is from the correct System and Component
				if frm.SystemID() != targetSystem || (targetComponent != 0 && frm.ComponentID() != targetComponent) {
					continue
				}

				// 2. Check for Mission Ack
				if msg, ok := frm.Message().(*common.MessageMissionAck); ok {
					// 3. Ensure the ACK is for our specific MISSION_TYPE_FENCE request
					if msg.MissionType == common.MAV_MISSION_TYPE_FENCE {
						if msg.Type == common.MAV_MISSION_ACCEPTED {
							log.Println("Geofence successfully cleared.")
							return true
						}
						log.Printf("Geofence clear failed. Result code: %v", msg.Type)
						return false
					}
				}
			}
		}
	}
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
func (s *DroneAPI) UploadGeofence(ctx context.Context, newFence []GeoFence, targetSystem uint8, targetComponent uint8) (bool, error) {
	log.Printf("Starting fence upload (%d points)...", len(newFence))

	// 1. Send the count to initiate the transaction
	s.drone.node.WriteMessageAll(&common.MessageMissionCount{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Count:           uint16(len(newFence)),
		MissionType:     common.MAV_MISSION_TYPE_FENCE,
	})

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return false, errors.New("node events channel closed")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok || frm.SystemID() != targetSystem {
				continue // Ignore messages from other systems
			}

			switch msg := frm.Message().(type) {

			case *common.MessageMissionRequestInt:
				if msg.MissionType != common.MAV_MISSION_TYPE_FENCE {
					continue
				}
				if int(msg.Seq) >= len(newFence) {
					return false, fmt.Errorf("drone requested out-of-bounds index: %d", msg.Seq)
				}

				// 2. Upload the specific requested point
				log.Printf("Uploading point #%d", msg.Seq)
				s.drone.node.WriteMessageAll(&common.MessageMissionItemInt{
					TargetSystem:    targetSystem,
					TargetComponent: targetComponent,
					Seq:             msg.Seq,
					Frame:           common.MAV_FRAME_GLOBAL,
					Command:         common.MAV_CMD_NAV_FENCE_POLYGON_VERTEX_INCLUSION,
					Param1:          float32(len(newFence)), // Total vertices
					X:               newFence[msg.Seq].Lat,
					Y:               newFence[msg.Seq].Lon,
					MissionType:     common.MAV_MISSION_TYPE_FENCE,
				})

			case *common.MessageMissionAck:
				if msg.MissionType != common.MAV_MISSION_TYPE_FENCE {
					continue
				}
				if msg.Type == common.MAV_MISSION_ACCEPTED {
					log.Println("Geofence upload complete and accepted.")
					return true, nil
				}
				return false, fmt.Errorf("mission rejected with code: %v", msg.Type)
			}
		}
	}
}
