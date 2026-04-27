package gomavlinkdroneapi

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

func (s *DroneAPI) UploadMission(droneChannel *gomavlib.Channel, mission []ardupilotmega.MessageMissionItemInt) bool {
	var rtn bool
	targetSysId := uint8(1)
	targetCompId := uint8(1)

	// Step 1: Tell the drone how many items we are uploading
	log.Println("Sending mission count...")
	s.drone.node.WriteMessageTo(droneChannel, &ardupilotmega.MessageMissionCount{
		TargetSystem:    targetSysId,
		TargetComponent: targetCompId,
		Count:           uint16(len(mission)),
		MissionType:     ardupilotmega.MAV_MISSION_TYPE_MISSION,
	})

	// Step 2 & 3: Listen for requests and answer them
	for evt := range s.drone.node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {

			switch msg := frm.Message().(type) {

			case *ardupilotmega.MessageMissionRequest:
				seq := msg.Seq
				if int(seq) < len(mission) {
					log.Printf("Uploading item %d...\n", seq)

					itemToSend := mission[seq]
					itemToSend.TargetSystem = targetSysId
					itemToSend.TargetComponent = targetCompId

					// Respond directly back to the specific channel requesting it
					s.drone.node.WriteMessageTo(frm.Channel, &itemToSend) // Use frm.Channel
				}

			case *ardupilotmega.MessageMissionAck:
				if msg.Type == ardupilotmega.MAV_MISSION_ACCEPTED {
					log.Println("Mission successfully uploaded and accepted!")
					rtn = true
				} else {
					log.Printf("Mission rejected with error code: %d\n", msg.Type)
				}

			}
		}
	}
	return rtn
}

func (s *DroneAPI) GetDroneChannel() *gomavlib.Channel {
	// Wait for the drone to send a heartbeat so we can capture its Channel
	log.Println("Waiting for heartbeat...")
	var droneChannel *gomavlib.Channel

	for evt := range s.drone.node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			if _, isHeartbeat := frm.Message().(*ardupilotmega.MessageHeartbeat); isHeartbeat {
				droneChannel = frm.Channel // Captured the actual communication channel!
				break
			}
		}
	}
	return droneChannel
}
