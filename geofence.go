package gomavlinkdroneapi

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
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
func (s *DroneAPI) UploadGeofence(ctx context.Context, fenceItems []ardupilotmega.MessageMissionItemInt, targetSystem uint8, targetComponent uint8) (bool, error) {
	log.Printf("Starting fence upload (%d points)...", len(fenceItems))
	// CRITICAL: Ensure Seq 0 is marked as 'Current' for ArduPilot's state machine
	if len(fenceItems) > 0 {
		fenceItems[0].Current = 1
	}

	// 1. Initiate transaction using ardupilotmega type
	s.drone.node.WriteMessageAll(&ardupilotmega.MessageMissionCount{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Count:           uint16(len(fenceItems)),
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
				continue
			}

			switch msg := frm.Message().(type) {

			// 2. Catch the ardupilotmega request (this fixes the "Cancelled" timeout)
			case *ardupilotmega.MessageMissionRequestInt:
				if msg.MissionType != common.MAV_MISSION_TYPE_FENCE {
					continue
				}

				seq := int(msg.Seq)
				if seq >= len(fenceItems) {
					return false, fmt.Errorf("out-of-bounds: %d", seq)
				}

				// 3. Send the item from your ardupilotmega slice
				item := fenceItems[seq]
				item.TargetSystem = targetSystem
				item.TargetComponent = targetComponent
				// MissionType is already set in your slice but we'll be safe:
				item.MissionType = common.MAV_MISSION_TYPE_FENCE

				log.Printf("Uploading point #%d", item.Seq)
				s.drone.node.WriteMessageAll(&item)

			// 4. Catch the ardupilotmega acknowledgement
			case *ardupilotmega.MessageMissionAck:
				if msg.MissionType != common.MAV_MISSION_TYPE_FENCE {
					continue
				}
				if msg.Type == common.MAV_MISSION_ACCEPTED {
					log.Println("Geofence accepted.")
					return true, nil
				}
				return false, fmt.Errorf("rejected with code: %v", msg.Type)
			}
		}
	}
}

func (s *DroneAPI) EnableGeofence(targetSystem uint8, targetComponent uint8) error {
	return s.drone.node.WriteMessageAll(&ardupilotmega.MessageCommandLong{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Command:         common.MAV_CMD_DO_FENCE_ENABLE,
		Param1:          1, // 1 to Enable, 0 to Disable
		Confirmation:    1,
	})
}
